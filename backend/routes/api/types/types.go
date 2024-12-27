package types

type ErrorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type AuthRequest struct {
	Username string `json:"username"`
}

type AuthResponse struct {
	Token string `json:"token,omitempty"`
	ErrorResponse
}

type Credentials struct {
	Token string `json:"token"`
}

type Event struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type DashboardRequest struct {
	Credentials
}

type DashboardResponse struct {
	Events []Event `json:"events"`
	ErrorResponse
}
