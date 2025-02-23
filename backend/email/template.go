package email

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"events/backend/database"
	"fmt"
	"log"
	"net/mail"
	"regexp"
	"strings"
	"text/template"

	"github.com/jhillyerd/enmime"
	"github.com/yuin/goldmark"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type Address struct {
	Name    string
	Address string
}

func ParseAddress(address string) (*Address, error) {
	parsed, err := mail.ParseAddress(address)
	if err != nil {
		return nil, fmt.Errorf("parse address: %v", err)
	}

	return &Address{Name: parsed.Name, Address: parsed.Address}, nil
}

type Email struct {
	To      *Address
	Cc      *Address
	Bcc     *Address
	ReplyTo *Address
	Subject string
	Body    string
}

func NewEmail(text string, keyValues map[string]string) (*Email, error) {
	text, err := renderPlain(text, keyValues)
	if err != nil {
		return nil, fmt.Errorf("render body: %v", err)
	}

	var e Email
	if err := e.parse(text); err != nil {
		return nil, err
	}

	return &e, nil
}

func NewEmailFromDB(ctx context.Context, template_id int64, participant_id int64) (*Email, error) {
	template, err := database.Default().Template(ctx, template_id)
	if err != nil {
		return nil, fmt.Errorf("get template: %v", err)
	}

	participant, err := database.Default().Participant(ctx, participant_id)
	if err != nil {
		return nil, err
	}

	var keyValues map[string]string
	err = json.Unmarshal([]byte(participant.Json.(string)), &keyValues)
	if err != nil {
		return nil, fmt.Errorf("unmarshal participant: %v", err)
	}

	return NewEmail(template.Body, keyValues)
}

func (e *Email) BodyHTML() string {
	return markdownToHTML(e.Body)
}

func (e *Email) SendGoogle(ctx context.Context, token string) error {
	// Build the message
	builder := enmime.Builder().From("Google ignores this and replaces it with the sender's address", "no-reply")
	if e.Subject != "" {
		builder = builder.Subject(e.Subject)
	}
	if e.To != nil {
		builder = builder.To(e.To.Name, e.To.Address)
	}
	if e.Cc != nil {
		builder = builder.CC(e.Cc.Name, e.Cc.Address)
	}
	if e.ReplyTo != nil {
		builder = builder.ReplyTo(e.ReplyTo.Name, e.ReplyTo.Address)
	}
	if e.Bcc != nil {
		builder = builder.BCC(e.Bcc.Name, e.Bcc.Address)
	}
	msg, err := builder.
		HTML([]byte(e.BodyHTML())).
		Text([]byte(e.Body)).
		Build()
	if err != nil {
		return fmt.Errorf("build message: %v", err)
	}

	var raw bytes.Buffer
	err = msg.Encode(&raw)
	if err != nil {
		return fmt.Errorf("encoding: %v", err)
	}

	log.Printf("raw: %s", raw.String())

	// Send the message
	gmailService, err := gmail.NewService(ctx, option.WithTokenSource(oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})))
	if err != nil {
		return fmt.Errorf("new gmail service: %v", err)
	}
	message := &gmail.Message{
		Raw: base64.URLEncoding.EncodeToString(raw.Bytes()),
	}
	_, err = gmailService.Users.Messages.Send("me", message).Do()
	if err != nil {
		return fmt.Errorf("send message: %v", err)
	}

	return nil
}

func (e *Email) parse(text string) error {
	lines := strings.Split(text, "\n")
	if lines[0] != "---" {
		return fmt.Errorf("missing '---' at the beginning of the email marking header start")
	}

	var err error
	for i, line := range lines[1:] {

		line = strings.TrimSpace(line)
		lowLine := strings.ToLower(line)

		if strings.HasPrefix(lowLine, "to:") {
			e.To, err = ParseAddress(strings.TrimSpace(line[3:]))
			if err != nil {
				return fmt.Errorf("parse to: %v", err)
			}
		} else if strings.HasPrefix(lowLine, "cc:") {
			e.Cc, err = ParseAddress(strings.TrimSpace(line[3:]))
			if err != nil {
				return fmt.Errorf("parse cc: %v", err)
			}
		} else if strings.HasPrefix(lowLine, "bcc:") {
			e.Bcc, err = ParseAddress(strings.TrimSpace(line[4:]))
			if err != nil {
				return fmt.Errorf("parse bcc: %v", err)
			}
		} else if strings.HasPrefix(lowLine, "reply-to:") {
			e.ReplyTo, err = ParseAddress(strings.TrimSpace(line[9:]))
			if err != nil {
				return fmt.Errorf("parse reply-to: %v", err)
			}
		} else if strings.HasPrefix(lowLine, "subject:") {
			e.Subject = strings.TrimSpace(line[8:])
		} else if line == "---" {
			e.Body = strings.Join(lines[i+2:], "\n")
			return nil
		} else {
			return fmt.Errorf("unknown header: %s", line)
		}
	}

	return fmt.Errorf("missing '---' at the end of the email marking header end")
}

func applyGenderTransformations(text string, keyValueData map[string]string) string {
	gender, ok := keyValueData[genderColumn]
	if !ok {
		return text
	}

	var pattern *regexp.Regexp
	if gender == genderMale {
		pattern = applyGenderMale
	} else if gender == genderFemale {
		pattern = applyGenderFemale
	} else {
		return text
	}

	return pattern.ReplaceAllString(text, "$1")
}

func renderPlain(input string, keyValueData map[string]string) (string, error) {
	// Parse the template
	tmpl, err := template.New("email").Parse(input)
	if err != nil {
		return "", fmt.Errorf("template parsing error: %v", err)
	}

	// Execute template with given data
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, keyValueData)
	if err != nil {
		return "", fmt.Errorf("template execution error: %v", err)
	}

	// Apply gender transformations
	out := applyGenderTransformations(buf.String(), keyValueData)
	return out, nil
}

func markdownToHTML(input string) string {
	var result bytes.Buffer
	goldmark.Convert([]byte(input), &result)
	return result.String()
}
