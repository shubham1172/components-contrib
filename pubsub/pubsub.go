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
	"strconv"

	"golang.org/x/sync/errgroup"
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

// BatchPublish publishes a batch of messages individually, as parallel requests.
// This implementation is used when the broker does not support batching.
// If a publish message fails, the whole batch will be failed.
func (p *DefaultMultiPubSub) BatchPublish(req *BatchPublishRequest) error {
	errs, _ := errgroup.WithContext(context.TODO())
	for i := range req.Data {
		id := i
		errs.Go(func() error {
			req := &PublishRequest{
				Data:        req.Data[id],
				PubsubName:  req.PubsubName,
				Topic:       req.Topic,
				Metadata:    req.Metadata,
				ContentType: req.ContentType,
			}
			return p.p.Publish(req)
		})
	}

	return errs.Wait()
}

// BulkSubsribe subscribes to a topic using a multi-message handler,
// that can be used to receive multiple messages at once.
func (p *DefaultMultiPubSub) BulkSubscribe(ctx context.Context, req SubscribeRequest, handler MultiMessageHandler) error {
	// Create a buffered channel to receive messages.
	msgs := make(chan *NewMessage, 100)

	opts := bulkMessageOptions{
		maxBatchCount:     getIntOrDefault(req.Metadata, maxBatchCountKey, 100),
		maxBatchSizeBytes: getIntOrDefault(req.Metadata, maxBatchSizeKey, 1024),
		maxBatchDelayMs:   getIntOrDefault(req.Metadata, maxBatchDelayMsKey, 5*1000),
	}

	p.p.Subscribe(ctx, req, func(ctx context.Context, msg *NewMessage) error {
		msgs <- msg
		return nil
	})

	// Start a goroutine to handle the messages.
	go processBulkMessages(ctx, msgs, opts, handler)
	return nil
}

// getIntOrDefault returns the value of the property with the given key as integer,
// or the default value if the property is not set.
func getIntOrDefault(m map[string]string, key string, def int) int {
	if val, ok := m[key]; ok {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return def
}
