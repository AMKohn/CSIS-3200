package processor

import (
	Mock "csis3200/mock/public"
	"time"
)

// Load mock data
var messagesDb = Mock.GetInitRequests(30)

/**
 * Gets the last 30 minutes of messages from the DB
 */
func GetRecentData() []map[string]interface{} {
	timeStamp := time.Now().Add(time.Duration(-30) * time.Minute).UnixNano() / 1000000

	// Find the first message that's new enough to include
	// The slice is cleaned up to remove messages older than 30 minutes and sorted
	// oldest first, so this is the most performant order
	i := 0
	mLen := len(messagesDb)

	// Check the next message since i is incremented afterward
	for i < mLen - 1 && messagesDb[i + 1] != nil && messagesDb[i]["timestamp"].(int64) <= timeStamp {
		i++
	}

	// Return this message and anything after (newer than) it
	// This is not a copy, but a collection of references
	return messagesDb[i:]
}

/**
 * Saves a new incoming log message into the DB
 */
func HandleMessage(message map[string]interface{}) {
	removeTimestamp := time.Now().Add(time.Duration(-30) * time.Minute).UnixNano() / 1000000

	// Update the message timestamp
	timeStamp := time.Now().UnixNano() / 1000000
	message["timestamp"] = timeStamp

	// Get rid of the oldest message if it's older than one hour
	// This works one message at a time for simplicity and performance. Removing more than one would
	// result in additional garbage collection and reduced performance without saving meaningful memory
	// since the space is already allocated and will likely be used soon
	// However, this means that the messagesDb slice will always be at it's maximum-yet size and never shrink
	if len(messagesDb) > 0 && messagesDb[0] != nil && messagesDb[0]["timestamp"].(int64) <= removeTimestamp {
		messagesDb = append(messagesDb[1:], message)
	} else {
		messagesDb = append(messagesDb, message)
	}
}
