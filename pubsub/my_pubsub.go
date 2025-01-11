//go:build !solution

package pubsub

import (
	"sync"
)

type SubscriptionImpl struct {
	service   *PubSubService
	topic     string
	callback  MsgHandler
	messageCh chan interface{}
	waitGroup *sync.WaitGroup
}

func (s *SubscriptionImpl) Unsubscribe() {
	s.service.removeSubscription(s.topic, s)
	close(s.messageCh)
}

func (s *SubscriptionImpl) listen() {
	go s.handleIncomingMessages()
}

func (s *SubscriptionImpl) handleIncomingMessages() {
	for msg := range s.messageCh {
		s.callback(msg)
		s.waitGroup.Done()
	}
}
