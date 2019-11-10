package processor

import (
	"csis3200/internal/app/dashboard"
	"time"
)

// Have a method to get JSON from the ingester
// "Save" log message in an array or something
// Send a copy of the message to the dashboard

// This shouldn't be an array, it needs to be something like a queue that only keeps up to 1000 messages, FIFO
const queueSize int = 1000

var messagesDb [queueSize]map[string]interface{}
var currentIndex = 0

func SaveMessage(message map[string]interface{}) {
	timeStamp := time.Now().UnixNano() / 1000000
	message["timestamp"] = timeStamp

	messagesDb[currentIndex] = message

	currentIndex++

	if currentIndex >= queueSize {
		currentIndex = 0
	}
}

func HandleMessage(message map[string]interface{}) {
	// "Save" the message
	SaveMessage(message)
	// Send a copy to the dashboard
	dashboard.HandleMessage(message)
}
