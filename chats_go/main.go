package main

import (
	"chats_go/api/controllers"
	"chats_go/application/services"
	"chats_go/infrastructure"
	"chats_go/infrastructure/jobs"
	"chats_go/persistence"
	"fmt"
	"github.com/gin-gonic/gin"
	"sync"
	"time"
)

var (
	backgroundJobMutex sync.Mutex
)

func main() {
	infrastructure.LoadEnvironmentVariables()
	db := persistence.InitDB()
	persistence.Migrate(db)
	redisClient := persistence.ConnectToRedis()
	elasticSearch := infrastructure.NewElasticSearch()
	elasticSearch.CreateIndex("messages")

	messageBroker := infrastructure.NewMessageBroker(db, elasticSearch)
	messageBroker.CreateQueueOnStartup("chats")
	messageBroker.CreateQueueOnStartup("messages")
	messageBroker.CreateQueueOnStartup("chats_count")
	go messageBroker.ConsumeFromChats()
	go messageBroker.ConsumeFromMessages()

	go runBackgroundJobs(services.NewChatService(db, redisClient, messageBroker))

	router := gin.Default()

	chatController := controllers.NewChatController(services.NewChatService(db, redisClient, messageBroker))
	router.POST("api/chats/create", chatController.Create)
	router.GET("api/chats/get", chatController.Get)

	messageController := controllers.NewMessageController(services.NewMessageService(db, redisClient, messageBroker, elasticSearch))
	router.POST("api/messages/create", messageController.Create)
	router.GET("api/messages/get", messageController.Get)
	router.POST("api/messages/search", messageController.Search)
	router.PUT("api/messages/update", messageController.Update)

	router.Run()
}

func runBackgroundJobs(chatService *services.ChatService) {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Check if a background job is already running
			if backgroundJobMutex.TryLock() {
				// If not, start a new job
				go func() {
					defer backgroundJobMutex.Unlock()
					jobs.CountChatsMessages(chatService)
					jobs.CountApplicationsChats(chatService)
				}()
			} else {
				fmt.Println("Previous background job is still running. Skipping...")
			}
		}
	}
}
