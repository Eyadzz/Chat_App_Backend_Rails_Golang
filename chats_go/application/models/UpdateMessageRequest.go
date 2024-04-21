package models

type UpdateMessageRequest struct {
	ApplicationToken string `json:"application_token" validate:"required"`
	ChatNumber       int    `json:"chat_number" validate:"required"`
	MessageNumber    int    `json:"message_number" validate:"required"`
	Message          string `json:"message" validate:"required"`
}
