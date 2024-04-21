package jobs

import (
	"chats_go/application/services"
	"log"
)

func CountChatsMessages(service *services.ChatService) {
	log.Printf("[Background Job] CountChatsMessages [START]")
	service.CountChatsMessages()
	log.Printf("[Background Job] CountChatsMessages [END]")
}

func CountApplicationsChats(service *services.ChatService) {
	log.Printf("[Background Job] CountApplicationsChats [START]")
	service.CountApplicationsChats()
	log.Printf("[Background Job] CountApplicationsChats [END]")
}
