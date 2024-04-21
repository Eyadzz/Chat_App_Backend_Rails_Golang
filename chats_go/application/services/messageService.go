package services

import (
	"chats_go/application/dtos"
	"chats_go/application/models"
	"chats_go/domain"
	"chats_go/infrastructure"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	"time"
)

type MessageService struct {
	DB            *gorm.DB
	RedisClient   *redis.Client
	MessageBroker *infrastructure.MessageBroker
	EsClient      *infrastructure.ElasticSearch
}

func NewMessageService(db *gorm.DB, redisClient *redis.Client, messageBroker *infrastructure.MessageBroker, esClient *infrastructure.ElasticSearch) *MessageService {
	return &MessageService{
		DB:            db,
		RedisClient:   redisClient,
		MessageBroker: messageBroker,
		EsClient:      esClient,
	}
}

func (MessageService MessageService) Add(request models.MessageRequest) (uint, error) {
	var chatID uint
	err := MessageService.DB.Model(&domain.Chat{}).
		Select("id").
		Where("application_token = ? AND number = ?", request.ApplicationToken, request.ChatNumber).
		First(&chatID).
		Error

	if err != nil {
		return 0, err
	}

	var number = uint(MessageService.RedisClient.Incr(context.Background(), fmt.Sprintf("chat_id:%s", chatID)).Val())

	message := domain.Message{ChatID: chatID, Content: request.Content, Number: number, CreatedAt: time.Now()}

	messageJson, _ := json.Marshal(message)
	go MessageService.MessageBroker.PublishToQueue("messages", messageJson)

	return number, nil
}

func (MessageService MessageService) Get(applicationToken string, chatNumber int) ([]dtos.MessageDto, error) {
	var chatID uint
	err := MessageService.DB.Model(&domain.Chat{}).
		Select("id").
		Where("application_token = ? AND number = ?", applicationToken, chatNumber).
		First(&chatID).
		Error

	if err != nil {
		return nil, err
	}

	var messages []domain.Message
	result := MessageService.DB.Find(&messages, "chat_id = ?", chatID)

	messageDtos := make([]dtos.MessageDto, len(messages))
	for i, message := range messages {
		messageDtos[i] = dtos.MessageDto{
			Number:    message.Number,
			Content:   message.Content,
			CreatedAt: message.CreatedAt.String(),
		}
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return messageDtos, nil
}

func (MessageService MessageService) Update(request models.UpdateMessageRequest) error {
	var chatID uint
	err := MessageService.DB.Model(&domain.Chat{}).
		Select("id").
		Where("application_token = ? AND number = ?", request.ApplicationToken, request.ChatNumber).
		First(&chatID).
		Error

	if err != nil {
		return err
	}

	var message domain.Message
	result := MessageService.DB.First(&message, "chat_id = ? AND number = ?", chatID, request.MessageNumber)

	if result.Error != nil {
		return result.Error
	}

	message.Content = request.Message
	result = MessageService.DB.Save(&message)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (MessageService MessageService) Search(request models.MessageRequest) ([]dtos.MessageDto, error) {
	var chatID uint
	err := MessageService.DB.Model(&domain.Chat{}).
		Select("id").
		Where("application_token = ? AND number = ?", request.ApplicationToken, request.ChatNumber).
		First(&chatID).
		Error

	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`{"query": {"bool": {"must": [{"term": {"chat_id": {"value": %d}}},{"match": {"message": "%s"}}]}}}`, chatID, request.Content)

	result := MessageService.EsClient.Search("messages", query)

	return convertEsResultToObject(result)
}

func convertEsResultToObject(esResponse *esapi.Response) ([]dtos.MessageDto, error) {
	// Extract the response body
	var responseBody map[string]interface{}
	if err := json.NewDecoder(esResponse.Body).Decode(&responseBody); err != nil {
		return nil, err
	}

	log.Printf("Elasticsearch response: %v", responseBody)

	// Extract hits from the Elasticsearch response
	hits, ok := responseBody["hits"].(map[string]interface{})["hits"].([]interface{})
	if !ok {
		log.Println("Hits not found or not in expected format")
		return nil, errors.New("unexpected format of Elasticsearch response")
	}

	log.Printf("Number of hits: %d", len(hits))

	// Extract message DTOs from the hits
	var messages []dtos.MessageDto
	for _, hit := range hits {
		hitData, ok := hit.(map[string]interface{})
		if !ok {
			continue
		}

		source, ok := hitData["_source"].(map[string]interface{})
		if !ok {
			log.Println("_source field not found or not in expected format")
			continue
		}

		log.Printf("source: %v", source)
		message := dtos.MessageDto{
			Number:    uint(source["number"].(float64)),
			Content:   source["message"].(string),
			CreatedAt: source["created_at"].(string),
		}
		messages = append(messages, message)
	}

	return messages, nil
}
