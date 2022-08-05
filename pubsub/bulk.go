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

const (
	// maxBatchSizeKey is the maximum number of messages to be published in a batch.
	maxBatchSizeKey = "maxBatchSize"
	// maxBatchTimeoutMsKey is the maximum time to wait before sending a batch of messages.
	// This is used to prevent waiting for individual messages for a large time.
	maxFetchWaitMsKey = "maxBatchTimeoutMs"
)
