package services

import (
	"chats_go/application/models"
	"chats_go/domain"
	"chats_go/infrastructure"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"time"
)

type ChatService struct {
	DB            *gorm.DB
	RedisClient   *redis.Client
	MessageBroker *infrastructure.MessageBroker
}

func NewChatService(db *gorm.DB, redisClient *redis.Client, messageBroker *infrastructure.MessageBroker) *ChatService {
	return &ChatService{
		DB:            db,
		RedisClient:   redisClient,
		MessageBroker: messageBroker,
	}
}

func (ChatService ChatService) Add(request models.ChatRequest) (uint, error) {
	var number = uint(ChatService.RedisClient.Incr(context.Background(), fmt.Sprintf("app_token:%s", request.ApplicationToken)).Val())

	chat := domain.Chat{ApplicationToken: request.ApplicationToken, Number: number, CreatedAt: time.Now()}
	chatJSON, _ := json.Marshal(chat)
	go infrastructure.MessageBroker{}.PublishToQueue("chats", chatJSON)

	return number, nil
}

func (ChatService ChatService) Get(applicationToken string) ([]uint, error) {
	var numbers []uint
	result := ChatService.DB.Model(&domain.Chat{}).
		Where("application_token = ?", applicationToken).
		Select("number").
		Find(&numbers)

	if result.Error != nil {
		return nil, result.Error
	}

	return numbers, nil
}

func (ChatService ChatService) CountApplicationsChats() {
	var chatCounts []models.ChatCount

	ChatService.DB.Model(&domain.Chat{}).
		Select("DISTINCT application_token, MAX(number) as number").
		Group("application_token").
		Scan(&chatCounts)

	for i := 0; i < len(chatCounts); i++ {
		chatJSON, _ := json.Marshal(chatCounts[i])
		go infrastructure.MessageBroker{}.PublishToQueue("chats_count", chatJSON)
	}
}

func (ChatService ChatService) CountChatsMessages() {
	var chats []domain.Chat
	ChatService.DB.Preload("Messages").Find(&chats)

	for _, chat := range chats {
		var maxMessage domain.Message
		ChatService.DB.Model(&chat).Order("number DESC").Limit(1).Association("Messages").Find(&maxMessage)

		chat.MessagesCount = maxMessage.Number

		ChatService.DB.Save(&chat)
	}
}
