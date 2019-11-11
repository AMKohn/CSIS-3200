package processor

import (
	"time"
)

// Have a method to get JSON from the ingester
// "Save" log message in an array or something
// Send a copy of the message to the dashboard

// This shouldn't be an array, it needs to be something like a queue that only keeps up to 1000 messages, FIFO
const queueSize int = 100000

var messagesDb [queueSize]map[string]interface{}
var currentIndex = 0

func GetRecentData() []map[string]interface{} {
	var messages []map[string]interface{}

	timeStamp := time.Now().Add(time.Duration(-30)*time.Minute).UnixNano() / 1000000

	for i := currentIndex; i < queueSize; i++ {
		if m := messages[i]; m["timestamp"].(int64) >= timeStamp {
			messages = append(messages, messagesDb[i])
		}
	}

	return messages
}

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
}
