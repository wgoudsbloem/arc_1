package pubsub

import (
	"errors"
)

// ErrSubscriberFailed is an error that can be handled
var ErrSubscriberFailed = errors.New("the subscriber could not be notfied")

type Subscriberer interface {
	Subscribe(fn Subscriber)
}

type Notifier interface {
	Notify(in interface{})
}

// PubSub is the interface for this package
type PubSub interface {
	Subscriberer
	Notifier
}

// Subscriber is the function that will be executed when Notify is called
type Subscriber func(in interface{}) error

// PubSuber is the default implementation of t PubSub
type PubSuber struct {
	subscribers []Subscriber
}

// Subscribe to the PubSub
func (ps *PubSuber) Subscribe(fn Subscriber) {
	ps.subscribers = append(ps.subscribers, fn)
}

// Notify to the PubSub
func (ps *PubSuber) Notify(in interface{}) {
	for i, subscriber := range ps.subscribers {
		if err := subscriber(in); err != nil {
			if err == ErrSubscriberFailed {
				//remove subscriber from list
				ps.subscribers = append(ps.subscribers[:i], ps.subscribers[i+1:]...)
			}
			//log.Errorf("subscriber failed %v", err)
		}
	}
}
