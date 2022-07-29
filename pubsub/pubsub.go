/*
Copyright 2021 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pubsub

import (
	"context"
)

// PubSub is the interface for message buses.
type PubSub interface {
	MultiPubSub
	Init(metadata Metadata) error
	Features() []Feature
	Publish(req *PublishRequest) error
	Subscribe(ctx context.Context, req SubscribeRequest, handler Handler) error
	Close() error
}

type MultiPubSub interface {
	BatchPublish(req *BatchPublishRequest) error
	BulkSubscribe(ctx context.Context, req SubscribeRequest, handler MultiMessageHandler) error
}

// Handler is the handler used to invoke the app handler.
type Handler func(ctx context.Context, msg *NewMessage) error

// MultiMessageHandler is the handler used to invoke the app handler.
type MultiMessageHandler func(ctx context.Context, msg []*NewMessage) error

type DefaultMultiPubSub struct {
	p PubSub
}

// NewDefaultBulkStore build a default bulk store.
func NewDefaultMultiPubsub(pubsub PubSub) DefaultMultiPubSub {
	defaultMultiPubSub := DefaultMultiPubSub{}
	defaultMultiPubSub.p = pubsub

	return defaultMultiPubSub
}

// BatchPublish publishes a batch of messages one by one.
// This implementation is used when the broker does not support batching.
// If a publish message fails, the whole batch will be failed.
func (p *DefaultMultiPubSub) BatchPublish(req *BatchPublishRequest) error {
	for _, msg := range req.Messages {
		req := &PublishRequest{
			Data:        msg.Data,
			PubsubName:  req.PubsubName,
			Topic:       msg.Topic,
			Metadata:    req.Metadata,
			ContentType: msg.ContentType,
		}
		err := p.p.Publish(req)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *DefaultMultiPubSub) BulkSubscribe(tx context.Context, req SubscribeRequest, handler MultiMessageHandler) error {
	return nil
}
