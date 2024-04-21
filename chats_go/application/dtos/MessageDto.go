package dtos

type MessageDto struct {
	Number    uint   `json:"message_number"`
	Content   string `json:"message"`
	CreatedAt string `json:"created_at"`
}
