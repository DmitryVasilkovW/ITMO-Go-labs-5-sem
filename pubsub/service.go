package pubsub

import (
	"context"
	"errors"
	"sync"
)

type PubSubService struct {
	topicSubscriptions map[string][]*SubscriptionImpl
	lock               sync.RWMutex
	waitGroup          sync.WaitGroup
	isClosed           bool
}

func NewPubSub() PubSub {
	return &PubSubService{
		topicSubscriptions: make(map[string][]*SubscriptionImpl),
	}
}

func (p *PubSubService) removeSubscription(topic string, subscription *SubscriptionImpl) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if subscribers, exists := p.topicSubscriptions[topic]; exists {
		p.subscribeHandler(topic, subscription, subscribers)
	}
}

func (p *PubSubService) subscribeHandler(topic string, subscription *SubscriptionImpl, subscribers []*SubscriptionImpl) {
	for i, s := range subscribers {
		if s == subscription {
			p.topicSubscriptions[topic] = append(subscribers[:i], subscribers[i+1:]...)
			break
		}
	}
}

func (p *PubSubService) Subscribe(topic string, callback MsgHandler) (Subscription, error) {
	if p.isServiceClosed() {
		return nil, errors.New("closed")
	}

	subscription := p.createSubscription(topic, callback)
	subscription.listen()
	p.registerSubscription(topic, subscription)

	return subscription, nil
}

func (p *PubSubService) registerSubscription(topic string, subscription *SubscriptionImpl) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.topicSubscriptions[topic] = append(p.topicSubscriptions[topic], subscription)
}

func (p *PubSubService) createSubscription(topic string, callback MsgHandler) *SubscriptionImpl {
	return &SubscriptionImpl{
		service:   p,
		topic:     topic,
		callback:  callback,
		messageCh: make(chan interface{}, 100),
		waitGroup: &p.waitGroup,
	}
}

func (p *PubSubService) Publish(topic string, message interface{}) error {
	if p.isServiceClosed() {
		return errors.New("closed")
	}

	p.sendMessageToTopicSubscribers(topic, message)
	return nil
}

func (p *PubSubService) sendMessageToTopicSubscribers(topic string, message interface{}) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	if subscribers, exists := p.topicSubscriptions[topic]; exists {
		p.sand(subscribers, message)
	}
}

func (p *PubSubService) sand(subscribers []*SubscriptionImpl, message interface{}) {
	for _, subscription := range subscribers {
		p.waitGroup.Add(1)
		subscription.messageCh <- message
	}
}

func (p *PubSubService) Close(ctx context.Context) error {
	if p.isServiceClosed() {
		return errors.New("closed")
	}

	p.terminateAllSubscriptions()
	return p.awaitPendingMessages(ctx)
}

func (p *PubSubService) awaitPendingMessages(ctx context.Context) error {
	completionSignal := make(chan struct{})
	go pending(p, completionSignal)

	select {
	case <-completionSignal:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func pending(p *PubSubService, completionSignal chan struct{}) {
	p.waitGroup.Wait()
	close(completionSignal)
}

func (p *PubSubService) terminateAllSubscriptions() {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.isClosed = true

	for _, subscribers := range p.topicSubscriptions {
		closeHandler(subscribers)
	}
}

func closeHandler(subscribers []*SubscriptionImpl) {
	for _, subscription := range subscribers {
		close(subscription.messageCh)
	}
}

func (p *PubSubService) isServiceClosed() bool {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.isClosed
}
