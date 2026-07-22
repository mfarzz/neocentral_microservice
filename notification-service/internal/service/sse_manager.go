package service

import (
	"log"
	"sync"
)

type SSEManager struct {
	clients map[string]map[chan []byte]bool
	mu      sync.RWMutex
}

func NewSSEManager() *SSEManager {
	return &SSEManager{
		clients: make(map[string]map[chan []byte]bool),
	}
}

func (m *SSEManager) AddClient(userID string, ch chan []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.clients[userID]; !ok {
		m.clients[userID] = make(map[chan []byte]bool)
	}
	m.clients[userID][ch] = true
	log.Printf("Client added for user %s. Total connections: %d", userID, len(m.clients[userID]))
}

func (m *SSEManager) RemoveClient(userID string, ch chan []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if userClients, ok := m.clients[userID]; ok {
		delete(userClients, ch)
		close(ch)
		if len(userClients) == 0 {
			delete(m.clients, userID)
		}
		log.Printf("Client removed for user %s", userID)
	}
}

func (m *SSEManager) SendToUser(userID string, payload []byte) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if userClients, ok := m.clients[userID]; ok {
		for ch := range userClients {
			// Non-blocking send to prevent hanging on slow clients
			select {
			case ch <- payload:
			default:
				log.Printf("Warning: Client for user %s is too slow, dropping message", userID)
			}
		}
	}
}
