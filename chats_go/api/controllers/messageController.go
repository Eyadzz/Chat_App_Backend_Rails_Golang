package controllers

import (
	"chats_go/application/models"
	"chats_go/application/services"
	"github.com/gin-gonic/gin"
	"strconv"
)

type MessageController struct {
	messageService *services.MessageService
}

func NewMessageController(messageService *services.MessageService) *MessageController {
	return &MessageController{messageService: messageService}
}

func (MessageController MessageController) Create(c *gin.Context) {
	var messageRequest models.MessageRequest
	if err := c.BindJSON(&messageRequest); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	messageNumber, err := MessageController.messageService.Add(messageRequest)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": messageNumber})
}

func (MessageController MessageController) Get(c *gin.Context) {

	chatNumber, _ := strconv.Atoi(c.Query("chat_number"))
	messages, err := MessageController.messageService.Get(c.Query("application_token"), chatNumber)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": messages})
}

func (MessageController MessageController) Update(c *gin.Context) {
	var messageRequest models.UpdateMessageRequest
	if err := c.BindJSON(&messageRequest); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	err := MessageController.messageService.Update(messageRequest)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": "Message updated"})
}

func (MessageController MessageController) Search(c *gin.Context) {
	var messageRequest models.MessageRequest
	if err := c.BindJSON(&messageRequest); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	results, err := MessageController.messageService.Search(messageRequest)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": results})
}
