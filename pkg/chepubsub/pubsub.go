// Package chepubsub provides a simple in-memory publish-subscribe event system.
package chepubsub

import (
	"context"
	"sync"
)

// Message represents a pub/sub message.
type Message struct {
	Topic string
	Data  interface{}
}

// Handler is a function that handles messages.
type Handler func(msg Message)

// Subscriber represents a subscription.
type Subscriber struct {
	id      string
	topic   string
	handler Handler
	ch      chan Message
	done    chan struct{}
}

// PubSub is an in-memory publish-subscribe system.
type PubSub struct {
	mu          sync.RWMutex
	subscribers map[string]map[string]*Subscriber // topic -> id -> subscriber
	bufferSize  int
	nextID      int
}

// New creates a new PubSub instance.
// bufferSize is the size of each subscriber's message buffer.
func New(bufferSize int) *PubSub {
	return &PubSub{
		subscribers: make(map[string]map[string]*Subscriber),
		bufferSize:  bufferSize,
	}
}

// Publish publishes a message to a topic.
// Returns the number of subscribers that received the message.
func (ps *PubSub) Publish(topic string, data interface{}) int {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	subs, ok := ps.subscribers[topic]
	if !ok {
		return 0
	}

	msg := Message{
		Topic: topic,
		Data:  data,
	}

	count := 0
	for _, sub := range subs {
		select {
		case sub.ch <- msg:
			count++
		default:
			// Subscriber buffer full, skip
		}
	}

	return count
}

// PublishSync publishes a message synchronously to all subscribers.
// Blocks until all subscribers have processed the message or context is cancelled.
func (ps *PubSub) PublishSync(ctx context.Context, topic string, data interface{}) error {
	ps.mu.RLock()
	subs, ok := ps.subscribers[topic]
	if !ok {
		ps.mu.RUnlock()
		return nil
	}

	// Copy subscribers to avoid holding lock during synchronous calls
	handlers := make([]Handler, 0, len(subs))
	for _, sub := range subs {
		handlers = append(handlers, sub.handler)
	}
	ps.mu.RUnlock()

	msg := Message{
		Topic: topic,
		Data:  data,
	}

	// Call handlers synchronously
	for _, handler := range handlers {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Run handler in goroutine to detect context cancellation
		done := make(chan struct{})
		go func() {
			handler(msg)
			close(done)
		}()

		select {
		case <-done:
			// Handler completed successfully
		case <-ctx.Done():
			// Context cancelled while handler was running
			return ctx.Err()
		}
	}

	return nil
}

// Subscribe subscribes to a topic with a handler function.
// Returns a Subscriber that can be used to unsubscribe.
func (ps *PubSub) Subscribe(topic string, handler Handler) *Subscriber {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if ps.subscribers[topic] == nil {
		ps.subscribers[topic] = make(map[string]*Subscriber)
	}

	ps.nextID++
	id := string(rune(ps.nextID))

	sub := &Subscriber{
		id:      id,
		topic:   topic,
		handler: handler,
		ch:      make(chan Message, ps.bufferSize),
		done:    make(chan struct{}),
	}

	ps.subscribers[topic][id] = sub

	// Start handler goroutine
	go sub.run()

	return sub
}

// Unsubscribe removes a subscriber.
func (ps *PubSub) Unsubscribe(sub *Subscriber) {
	if sub == nil {
		return
	}

	ps.mu.Lock()
	defer ps.mu.Unlock()

	if subs, ok := ps.subscribers[sub.topic]; ok {
		if _, exists := subs[sub.id]; exists {
			delete(subs, sub.id)
			if len(subs) == 0 {
				delete(ps.subscribers, sub.topic)
			}
			close(sub.done)
		}
	}
}

// UnsubscribeAll removes all subscribers from a topic.
func (ps *PubSub) UnsubscribeAll(topic string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if subs, ok := ps.subscribers[topic]; ok {
		for _, sub := range subs {
			close(sub.done)
		}
		delete(ps.subscribers, topic)
	}
}

// Topics returns all topics that have active subscribers.
func (ps *PubSub) Topics() []string {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	topics := make([]string, 0, len(ps.subscribers))
	for topic := range ps.subscribers {
		topics = append(topics, topic)
	}
	return topics
}

// SubscriberCount returns the number of subscribers for a topic.
func (ps *PubSub) SubscriberCount(topic string) int {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if subs, ok := ps.subscribers[topic]; ok {
		return len(subs)
	}
	return 0
}

// Close closes all subscriptions and cleans up.
func (ps *PubSub) Close() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for topic, subs := range ps.subscribers {
		for _, sub := range subs {
			close(sub.done)
		}
		delete(ps.subscribers, topic)
	}
}

func (s *Subscriber) run() {
	for {
		select {
		case msg := <-s.ch:
			s.handler(msg)
		case <-s.done:
			return
		}
	}
}

// Topic returns the topic this subscriber is subscribed to.
func (s *Subscriber) Topic() string {
	return s.topic
}
