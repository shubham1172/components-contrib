/*
Copyright 2022 The Dapr Authors
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
	"time"
)

const (
	// batchPublishKeepOrderKey is used to indicate whether to keep the order of messages in a batch publish request.
	batchPublishKeepOrderKey = "batchPublishKeepOrder"
	// maxBatchCountKey is the maximum number of messages to be published in a batch.
	maxBatchCountKey = "maxBatchCount"
	// maxBatchSizeBytesKey is the maximum size of a batch in bytes.
	maxBatchSizeKey = "maxBatchSizeBytes"
	// maxBatchDelayMsKey is the maximum delay in milliseconds before publishing a batch.
	maxBatchDelayMsKey = "maxBatchDelayMs"
)

type bulkMessageOptions struct {
	maxBatchCount     int
	maxBatchSizeBytes int
	maxBatchDelayMs   int
}

// writes messages to the MultiMessageHandler.
func flushMessages(messages []*NewMessage, handler MultiMessageHandler) {
	if len(messages) > 0 {
		// TODO: log error
		handler(context.Background(), messages)
	}
}

// processBulkMessages reads messages from the channel and publishes them to MultiMessageHandler.
// It buffers messages in memory and publishes them in batches.
// TODO: Do not just use maxBatchCount, but also introduce maxBatchSizeBytes and maxBatchTimeoutMs.
func processBulkMessages(ctx context.Context, msgChan <-chan *NewMessage, opts bulkMessageOptions, handler MultiMessageHandler) {
	var messages []*NewMessage
	var currSize int

	ticker := time.NewTicker(time.Duration(opts.maxBatchDelayMs) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			flushMessages(messages, handler)
			return
		case msg := <-msgChan:
			currSize += 0 // TODO: figure out the size of message
			messages = append(messages, msg)
			if len(messages) >= opts.maxBatchCount || currSize >= opts.maxBatchSizeBytes {
				flushMessages(messages, handler)
			}
		case <-ticker.C:
			flushMessages(messages, handler)
		}
	}
}
