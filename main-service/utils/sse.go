package utils

import (
	"bufio"
	"fmt"
	"sync"
)

var sseClients = make(map[string]chan string)
var mu sync.RWMutex

// RegisterSSE creates a channel for a given search_id
func RegisterSSE(searchID string) {
	mu.Lock()
	defer mu.Unlock()
	sseClients[searchID] = make(chan string, 10)
}

// GetSSEChannel retrieves channel for given search_id
func GetSSEChannel(searchID string) (chan string, bool) {
	mu.RLock()
	defer mu.RUnlock()
	ch, ok := sseClients[searchID]
	return ch, ok
}

// RemoveSSE cleans up the channel
func RemoveSSE(searchID string) {
	mu.Lock()
	defer mu.Unlock()
	if ch, ok := sseClients[searchID]; ok {
		close(ch)
		delete(sseClients, searchID)
	}
}

// WriteSSE writes data to bufio.Writer for streaming
func WriteSSE(w *bufio.Writer, ch chan string) {
	for msg := range ch {
		fmt.Fprintf(w, "data: %s\n\n", msg)
		w.Flush()
	}
}
