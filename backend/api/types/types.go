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

type Project struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type DashboardResponse struct {
	Projects []Project `json:"projects"`
}
