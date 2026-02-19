package email

import (
	"context"
	"log"
	"sync"
)

// MockSender implements a mock email sender for testing.
type MockSender struct {
	mu       sync.Mutex
	messages []*Message
}

// NewMockSender creates a new mock email sender.
func NewMockSender() *MockSender {
	return &MockSender{
		messages: make([]*Message, 0),
	}
}

// Send logs the email and stores it for inspection.
func (s *MockSender) Send(ctx context.Context, msg *Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	log.Printf("[MockEmail] To: %s, Subject: %s", msg.To, msg.Subject)
	if msg.TextBody != "" {
		log.Printf("[MockEmail] Body:\n%s", msg.TextBody)
	}
	s.messages = append(s.messages, msg)

	return nil
}

// SendBulk logs multiple emails and stores them.
func (s *MockSender) SendBulk(ctx context.Context, msgs []*Message) error {
	for _, msg := range msgs {
		if err := s.Send(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}

// Messages returns all sent messages.
func (s *MockSender) Messages() []*Message {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]*Message, len(s.messages))
	copy(result, s.messages)
	return result
}

// Clear removes all stored messages.
func (s *MockSender) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.messages = make([]*Message, 0)
}

// LastMessage returns the most recently sent message.
func (s *MockSender) LastMessage() *Message {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.messages) == 0 {
		return nil
	}
	return s.messages[len(s.messages)-1]
}
