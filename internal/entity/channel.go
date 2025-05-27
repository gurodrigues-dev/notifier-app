package entity

type Channel struct {
	ID       int    `json:"id"`
	Platform string `json:"platform" validate:"required"`
	TargetID string `json:"target_id" validate:"required"`
	Group    string `json:"group" validate:"required"`
}
