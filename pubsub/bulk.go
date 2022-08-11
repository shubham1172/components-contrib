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

import "context"

const (
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

// processBulkMessages reads messages from the channel and publishes them to MultiMessageHandler.
// It buffers messages in memory and publishes them in batches.
// TODO: Do not just use maxBatchCount, but also introduce maxBatchSizeBytes and maxBatchTimeoutMs.
func processBulkMessages(ctx context.Context, msgChan <-chan *NewMessage, opts bulkMessageOptions, handler MultiMessageHandler) {
	var messages []*NewMessage
	for msg := range msgChan {
		messages = append(messages, msg)
		if len(messages) == opts.maxBatchCount {
			_ = handler(ctx, messages)
			// TODO: Handle error.
			messages = nil
		}
	}
	// Handle remaining messages.
	if len(messages) > 0 {
		_ = handler(ctx, messages)
		// TODO: Handle error.
	}
}
