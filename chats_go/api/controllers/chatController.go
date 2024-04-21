package controllers

import (
	"chats_go/application/models"
	"chats_go/application/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ChatController struct {
	chatService *services.ChatService
}

func NewChatController(chatService *services.ChatService) *ChatController {
	return &ChatController{chatService: chatService}
}

func (ChatController ChatController) Create(c *gin.Context) {

	var createChatRequest models.ChatRequest
	c.Bind(&createChatRequest)

	number, err := ChatController.chatService.Add(createChatRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": number,
	})
}

func (ChatController ChatController) Get(c *gin.Context) {

	applicationToken := c.Query("application_token")

	numbers, err := ChatController.chatService.Get(applicationToken)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": numbers,
	})
}
