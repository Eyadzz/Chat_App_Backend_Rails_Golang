package models

type MessageRequest struct {
	ApplicationToken string `json:"application_token"`
	ChatNumber       uint   `json:"chat_number"`
	Content          string `json:"message"`
}
