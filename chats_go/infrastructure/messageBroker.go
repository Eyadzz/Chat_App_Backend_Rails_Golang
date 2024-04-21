package infrastructure

import (
	"chats_go/domain"
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

type MessageBroker struct {
	DB            *gorm.DB
	ElasticSearch *ElasticSearch
}

func NewMessageBroker(db *gorm.DB, es *ElasticSearch) *MessageBroker {
	return &MessageBroker{
		DB:            db,
		ElasticSearch: es,
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func (MessageBroker MessageBroker) CreateQueueOnStartup(name string) {

	connection, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer connection.Close()

	channel, err := connection.Channel()
	failOnError(err, "Failed to open a channel")
	defer channel.Close()

	_, err = channel.QueueDeclare(name, false, false, false, false, nil)
	failOnError(err, "Failed to declare chats queue")
}

func (MessageBroker MessageBroker) PublishToQueue(queueName string, message []byte) {
	connection, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer connection.Close()

	channel, err := connection.Channel()
	failOnError(err, "Failed to open a channel")
	defer channel.Close()

	queue, err := channel.QueueDeclare(queueName, false, false, false, false, nil)
	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = channel.PublishWithContext(ctx, "", queue.Name, false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		})
	failOnError(err, "Failed to publish a message")
}

func (MessageBroker MessageBroker) ConsumeFromChats() {
	connection, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer connection.Close()

	channel, err := connection.Channel()
	failOnError(err, "Failed to open a channel")
	defer channel.Close()

	messages, err := channel.Consume("chats", "", true, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")

	log.Printf("Consuming messages from queue: chats")

	var infinite chan struct{}

	go func() {
		for message := range messages {
			var chat domain.Chat
			json.Unmarshal(message.Body, &chat)
			log.Printf("Received chat: %+v", chat)

			MessageBroker.DB.Create(&chat)

			log.Printf("Chat saved to database: %+v", chat.ID)
		}
	}()

	<-infinite
}

func (MessageBroker MessageBroker) ConsumeFromMessages() {
	connection, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer connection.Close()

	channel, err := connection.Channel()
	failOnError(err, "Failed to open a channel")
	defer channel.Close()

	messages, err := channel.Consume("messages", "", true, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")

	log.Printf("Consuming messages from queue: messages")

	var infinite chan struct{}

	go func() {
		for message := range messages {
			var newMessage domain.Message
			json.Unmarshal(message.Body, &newMessage)

			MessageBroker.DB.Create(&newMessage)
			es_doc := fmt.Sprintf(`{"chat_id": %d, "number": %d, "message": "%s", "created_at":"%s"}`, newMessage.ChatID, newMessage.Number, newMessage.Content, newMessage.CreatedAt.String())
			MessageBroker.ElasticSearch.Insert("messages", es_doc)

			log.Printf("Message saved to database: %+v", newMessage.ID)
		}
	}()

	<-infinite
}
