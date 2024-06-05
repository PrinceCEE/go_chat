package models

type Auth struct {
	ModelMixin
	UserID   string `json:"user_id"`
	Password string `json:"password"`
}
