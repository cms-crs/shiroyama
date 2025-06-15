package dto

type CreateTeamRequest struct {
	TeamName    string
	Description string
	CreatedBy   string
}

type CreateTeamResponse struct {
	ID          string
	TeamName    string
	Description string
}
