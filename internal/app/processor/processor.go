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
	timeStamp := time.Now().Add(time.Duration(-30)*time.Minute).UnixNano() / 1000000

	i := currentIndex

	for i > 0 && messagesDb[i] != nil && messagesDb[i]["timestamp"].(int64) >= timeStamp {
		i--
	}

	// If the first item in the messagesDb is less than 30 minutes old, check and
	// see if we need to wrap around the end (the last item is recent enough)
	if i == 0 && messagesDb[queueSize-1] != nil && messagesDb[queueSize-1]["timestamp"].(int64) >= timeStamp {
		i := queueSize - 1

		for i >= currentIndex && messagesDb[i] != nil && messagesDb[i]["timestamp"].(int64) >= timeStamp {
			i--
		}
	}

	var messages []map[string]interface{}

	if i < currentIndex { // This is a normal slicing operation
		messages = messagesDb[i:currentIndex]
	} else { // Otherwise, i > currentIndex and we wrapped around the beginning, get the right sections
		messages = append(messagesDb[currentIndex:i], messagesDb[0:currentIndex]...)
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
