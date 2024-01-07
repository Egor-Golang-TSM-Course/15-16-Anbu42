package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	Email string = "email"
	SMS   string = "sms"
	Push  string = "push"
)

type UserPreferences struct {
	EmailEnabled bool
	SMSEnabled   bool
	PushEnabled  bool
}

type Notification struct {
	ctx   context.Context
	sms   chan string
	email chan string
	push  chan string
}

func (n *Notification) Send(notificType string, message string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-n.ctx.Done():
			return
		default:
			switch notificType {
			case SMS:
				n.sms <- message
			case Email:
				n.email <- message
			case Push:
				n.push <- message
			default:
				fmt.Println("undefined type")
			}
			return
		}
	}
}

func (n *Notification) Notify(ctx context.Context, message string, wg *sync.WaitGroup) {
	defer wg.Done()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Finishing")
			close(n.sms)
			close(n.email)
			close(n.push)
			os.Exit(0)
		case <-sigChan:
			cancel()
		case message, ok := <-n.email:
			if ok {
				fmt.Println("Received email notification:", message)
			}
		case message, ok := <-n.sms:
			if ok {
				fmt.Println("Received sms notification:", message)
			}
		case message, ok := <-n.push:
			if ok {
				fmt.Println("Received push notification:", message)
			}
		}
	}
}

func NewNotification(ctx context.Context) *Notification {
	return &Notification{
		ctx:   ctx,
		sms:   make(chan string),
		email: make(chan string),
		push:  make(chan string),
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	userPreferences := map[int]UserPreferences{
		1: {EmailEnabled: true, SMSEnabled: false, PushEnabled: true},
		2: {EmailEnabled: false, SMSEnabled: true, PushEnabled: false},
		3: {EmailEnabled: true, SMSEnabled: true, PushEnabled: true},
	}

	notification := NewNotification(ctx)

	// Start notification senders for each user based on preferences
	for userID, preferences := range userPreferences {
		wg.Add(1)
		go func(userID int, preferences UserPreferences) {
			defer wg.Done()
			message := fmt.Sprintf("Hello User %d! You have a new notification.", userID)

			if preferences.EmailEnabled {
				wg.Add(2)
				go notification.Send(Email, message, &wg)
				go notification.Notify(ctx, message, &wg)
			}

			if preferences.SMSEnabled {
				wg.Add(2)
				go notification.Send(SMS, message, &wg)
				go notification.Notify(ctx, message, &wg)
			}

			if preferences.PushEnabled {
				wg.Add(2)
				go notification.Send(Push, message, &wg)
				go notification.Notify(ctx, message, &wg)
			}
		}(userID, preferences)
	}

	// Simulate a change in user preferences or an external signal to stop notifications
	time.Sleep(5 * time.Second)
	fmt.Println("User preferences changed or external signal received. Stopping notifications...")
	cancel()

	wg.Wait()
}
