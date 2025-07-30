package utils

import (
	"bufio"
	"fmt"
	"sync"
)

var sseClients = make(map[string]chan string)
var mu sync.RWMutex

func RegisterSSE(searchID string) {
	mu.Lock()
	defer mu.Unlock()
	sseClients[searchID] = make(chan string, 10)
}

func GetSSEChannel(searchID string) (chan string, bool) {
	mu.RLock()
	defer mu.RUnlock()
	ch, ok := sseClients[searchID]
	return ch, ok
}

func RemoveSSE(searchID string) {
	mu.Lock()
	defer mu.Unlock()
	if ch, ok := sseClients[searchID]; ok {
		close(ch)
		delete(sseClients, searchID)
	}
}

func WriteSSE(w *bufio.Writer, ch chan string) {
	for msg := range ch {
		fmt.Fprintf(w, "data: %s\n\n", msg)
		w.Flush()
	}
}
