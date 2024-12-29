package types

type ErrorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type Event struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type EventResponse struct {
	Events []Event `json:"events"`
	ErrorResponse
}

type Profile struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	PictureURL string `json:"pictureURL"`
}
