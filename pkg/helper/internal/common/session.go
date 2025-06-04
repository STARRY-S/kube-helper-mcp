package common

import (
	"github.com/mark3labs/mcp-go/mcp"
)

// Implement your own ClientSession
type Session struct {
	id            string
	notifChannel  chan mcp.JSONRPCNotification
	isInitialized bool
	// Add custom fields for your application
}

func NewSession(id string) *Session {
	session := &Session{
		id:           id,
		notifChannel: make(chan mcp.JSONRPCNotification, 10),
	}
	return session
}

func (s *Session) SessionID() string {
	return s.id
}

func (s *Session) NotificationChannel() chan<- mcp.JSONRPCNotification {
	return s.notifChannel
}

func (s *Session) Initialize() {
	s.isInitialized = true
}

func (s *Session) Initialized() bool {
	return s.isInitialized
}
