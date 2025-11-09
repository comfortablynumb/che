package chepubsub

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/comfortablynumb/che/pkg/chetest"
)

func TestPubSub_Basic(t *testing.T) {
	ps := New(10)
	defer ps.Close()

	var received []string
	var mu sync.Mutex

	sub := ps.Subscribe("test", func(msg Message) {
		mu.Lock()
		defer mu.Unlock()
		received = append(received, msg.Data.(string))
	})
	defer ps.Unsubscribe(sub)

	count := ps.Publish("test", "hello")
	chetest.RequireEqual(t, count, 1)

	// Give handler time to process
	time.Sleep(10 * time.Millisecond)

	mu.Lock()
	chetest.RequireEqual(t, len(received), 1)
	chetest.RequireEqual(t, received[0], "hello")
	mu.Unlock()
}

func TestPubSub_MultipleSubscribers(t *testing.T) {
	ps := New(10)
	defer ps.Close()

	var received1, received2 []string
	var mu1, mu2 sync.Mutex

	sub1 := ps.Subscribe("test", func(msg Message) {
		mu1.Lock()
		defer mu1.Unlock()
		received1 = append(received1, msg.Data.(string))
	})
	defer ps.Unsubscribe(sub1)

	sub2 := ps.Subscribe("test", func(msg Message) {
		mu2.Lock()
		defer mu2.Unlock()
		received2 = append(received2, msg.Data.(string))
	})
	defer ps.Unsubscribe(sub2)

	count := ps.Publish("test", "hello")
	chetest.RequireEqual(t, count, 2)

	time.Sleep(10 * time.Millisecond)

	mu1.Lock()
	chetest.RequireEqual(t, len(received1), 1)
	chetest.RequireEqual(t, received1[0], "hello")
	mu1.Unlock()

	mu2.Lock()
	chetest.RequireEqual(t, len(received2), 1)
	chetest.RequireEqual(t, received2[0], "hello")
	mu2.Unlock()
}

func TestPubSub_MultipleTopics(t *testing.T) {
	ps := New(10)
	defer ps.Close()

	var topic1, topic2 []string
	var mu1, mu2 sync.Mutex

	sub1 := ps.Subscribe("topic1", func(msg Message) {
		mu1.Lock()
		defer mu1.Unlock()
		topic1 = append(topic1, msg.Data.(string))
	})
	defer ps.Unsubscribe(sub1)

	sub2 := ps.Subscribe("topic2", func(msg Message) {
		mu2.Lock()
		defer mu2.Unlock()
		topic2 = append(topic2, msg.Data.(string))
	})
	defer ps.Unsubscribe(sub2)

	ps.Publish("topic1", "message1")
	ps.Publish("topic2", "message2")

	time.Sleep(10 * time.Millisecond)

	mu1.Lock()
	chetest.RequireEqual(t, len(topic1), 1)
	chetest.RequireEqual(t, topic1[0], "message1")
	mu1.Unlock()

	mu2.Lock()
	chetest.RequireEqual(t, len(topic2), 1)
	chetest.RequireEqual(t, topic2[0], "message2")
	mu2.Unlock()
}

func TestPubSub_NoSubscribers(t *testing.T) {
	ps := New(10)
	defer ps.Close()

	count := ps.Publish("nonexistent", "hello")
	chetest.RequireEqual(t, count, 0)
}

func TestPubSub_Unsubscribe(t *testing.T) {
	ps := New(10)
	defer ps.Close()

	var received []string
	var mu sync.Mutex

	sub := ps.Subscribe("test", func(msg Message) {
		mu.Lock()
		defer mu.Unlock()
		received = append(received, msg.Data.(string))
	})

	ps.Publish("test", "first")
	time.Sleep(10 * time.Millisecond)

	ps.Unsubscribe(sub)

	ps.Publish("test", "second")
	time.Sleep(10 * time.Millisecond)

	mu.Lock()
	chetest.RequireEqual(t, len(received), 1)
	chetest.RequireEqual(t, received[0], "first")
	mu.Unlock()
}

func TestPubSub_UnsubscribeAll(t *testing.T) {
	ps := New(10)
	defer ps.Close()

	var count1, count2 int
	var mu sync.Mutex

	ps.Subscribe("test", func(msg Message) {
		mu.Lock()
		count1++
		mu.Unlock()
	})

	ps.Subscribe("test", func(msg Message) {
		mu.Lock()
		count2++
		mu.Unlock()
	})

	ps.Publish("test", "first")
	time.Sleep(10 * time.Millisecond)

	ps.UnsubscribeAll("test")

	ps.Publish("test", "second")
	time.Sleep(10 * time.Millisecond)

	mu.Lock()
	chetest.RequireEqual(t, count1, 1)
	chetest.RequireEqual(t, count2, 1)
	mu.Unlock()
}

func TestPubSub_Topics(t *testing.T) {
	ps := New(10)
	defer ps.Close()

	sub1 := ps.Subscribe("topic1", func(msg Message) {})
	defer ps.Unsubscribe(sub1)

	sub2 := ps.Subscribe("topic2", func(msg Message) {})
	defer ps.Unsubscribe(sub2)

	topics := ps.Topics()
	chetest.RequireEqual(t, len(topics), 2)

	// Topics can be in any order
	hasT1 := false
	hasT2 := false
	for _, t := range topics {
		if t == "topic1" {
			hasT1 = true
		}
		if t == "topic2" {
			hasT2 = true
		}
	}
	chetest.RequireEqual(t, hasT1, true)
	chetest.RequireEqual(t, hasT2, true)
}

func TestPubSub_SubscriberCount(t *testing.T) {
	ps := New(10)
	defer ps.Close()

	chetest.RequireEqual(t, ps.SubscriberCount("test"), 0)

	sub1 := ps.Subscribe("test", func(msg Message) {})
	defer ps.Unsubscribe(sub1)
	chetest.RequireEqual(t, ps.SubscriberCount("test"), 1)

	sub2 := ps.Subscribe("test", func(msg Message) {})
	defer ps.Unsubscribe(sub2)
	chetest.RequireEqual(t, ps.SubscriberCount("test"), 2)

	ps.Unsubscribe(sub1)
	chetest.RequireEqual(t, ps.SubscriberCount("test"), 1)
}

func TestPubSub_PublishSync(t *testing.T) {
	ps := New(10)
	defer ps.Close()

	var received []string
	var mu sync.Mutex

	ps.Subscribe("test", func(msg Message) {
		mu.Lock()
		defer mu.Unlock()
		received = append(received, msg.Data.(string))
	})

	ctx := context.Background()
	err := ps.PublishSync(ctx, "test", "hello")
	chetest.RequireEqual(t, err, nil)

	// Should be processed immediately in sync mode
	mu.Lock()
	chetest.RequireEqual(t, len(received), 1)
	chetest.RequireEqual(t, received[0], "hello")
	mu.Unlock()
}

func TestPubSub_PublishSyncWithContext(t *testing.T) {
	ps := New(10)
	defer ps.Close()

	var received int
	var mu sync.Mutex

	ps.Subscribe("test", func(msg Message) {
		mu.Lock()
		defer mu.Unlock()
		received++
		time.Sleep(100 * time.Millisecond)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := ps.PublishSync(ctx, "test", "hello")
	chetest.RequireEqual(t, err != nil, true)
}

func TestPubSub_BufferFull(t *testing.T) {
	ps := New(1) // Small buffer
	defer ps.Close()

	var wg sync.WaitGroup
	wg.Add(1)

	// Slow handler
	ps.Subscribe("test", func(msg Message) {
		time.Sleep(100 * time.Millisecond)
		wg.Done()
	})

	// First message should be buffered
	count := ps.Publish("test", "first")
	chetest.RequireEqual(t, count, 1)

	// Second message should fail due to full buffer
	count = ps.Publish("test", "second")
	chetest.RequireEqual(t, count, 0)

	wg.Wait()
}

func TestPubSub_Close(t *testing.T) {
	ps := New(10)

	sub1 := ps.Subscribe("topic1", func(msg Message) {})
	sub2 := ps.Subscribe("topic2", func(msg Message) {})

	chetest.RequireEqual(t, ps.SubscriberCount("topic1"), 1)
	chetest.RequireEqual(t, ps.SubscriberCount("topic2"), 1)

	ps.Close()

	chetest.RequireEqual(t, ps.SubscriberCount("topic1"), 0)
	chetest.RequireEqual(t, ps.SubscriberCount("topic2"), 0)
	chetest.RequireEqual(t, len(ps.Topics()), 0)

	// Verify subscribers are cleaned up
	chetest.RequireEqual(t, sub1 != nil, true)
	chetest.RequireEqual(t, sub2 != nil, true)
}

func TestSubscriber_Topic(t *testing.T) {
	ps := New(10)
	defer ps.Close()

	sub := ps.Subscribe("test-topic", func(msg Message) {})
	defer ps.Unsubscribe(sub)

	chetest.RequireEqual(t, sub.Topic(), "test-topic")
}

func TestPubSub_ConcurrentPublish(t *testing.T) {
	ps := New(100)
	defer ps.Close()

	var received int
	var mu sync.Mutex

	ps.Subscribe("test", func(msg Message) {
		mu.Lock()
		received++
		mu.Unlock()
	})

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ps.Publish("test", "message")
		}()
	}

	wg.Wait()
	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	chetest.RequireEqual(t, received, 100)
	mu.Unlock()
}

func TestPubSub_UnsubscribeNil(t *testing.T) {
	ps := New(10)
	defer ps.Close()

	// Should not panic
	ps.Unsubscribe(nil)
}
