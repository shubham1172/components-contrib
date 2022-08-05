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

// PublishRequest is the request to publish a message.
type PublishRequest struct {
	Data        []byte            `json:"data"`
	PubsubName  string            `json:"pubsubname"`
	Topic       string            `json:"topic"`
	Metadata    map[string]string `json:"metadata"`
	ContentType *string           `json:"contentType,omitempty"`
}

// BatchPublishMessage is a message from a BatchPublishRequest.
type BatchPublishMessage struct {
	Data []byte `json:"data"`
}

// BatchPublishRequest is the message to publish events data to pubsub topic
type BatchPublishRequest struct {
	PubsubName  string                `json:"pubsubname"`
	Topic       string                `json:"topic"`
	Messages    []BatchPublishMessage `json:"messages"`
	ContentType *string               `json:"contentType,omitempty"`
	Metadata    map[string]string     `json:"metadata"`
}

// SubscribeRequest is the request to subscribe to a topic.
type SubscribeRequest struct {
	Topic    string            `json:"topic"`
	Metadata map[string]string `json:"metadata"`
}

// BulkSubscribeRequest is the request to receive multiple messages.
type BulkSubscribeRequest struct {
	PubSubName string            `json:"pubsubname"`
	Metadata   map[string]string `json:"metadata"`
}

// NewMessage is an event arriving from a message bus instance.
type NewMessage struct {
	Data        []byte            `json:"data"`
	Topic       string            `json:"topic"`
	Metadata    map[string]string `json:"metadata"`
	ContentType *string           `json:"contentType,omitempty"`
}
