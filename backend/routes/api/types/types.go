package types

type ErrorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type Event struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Events []Event

type Profile struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	PictureURL string `json:"pictureURL"`
}

type Organizer struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Organizers []Organizer

type Participant struct {
	ID   string            `json:"id"`
	Data map[string]string `json:"data"`
}

type Participants struct {
	Participants []Participant `json:"participants"`
}

type Template struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Body string `json:"body"`
}

type Templates []Template
