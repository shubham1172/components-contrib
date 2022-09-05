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

// NewBatchPublishResponse creates a BatchPublishResponse with messages and the status.
// This is used when a broker does not support per-message status, so the status is applied to all messages.
func NewBatchPublishResponse(messages []BatchMessage, status BatchMessageResponseStatus) BatchPublishResponse {
	var responses []BatchMessageResponse

	for _, m := range messages {
		responses = append(responses, BatchMessageResponse{
			ID:     m.ID,
			Status: status,
		})
	}

	return BatchPublishResponse{
		Statuses: responses,
	}
}

// WithError adds an error to the BatchPublishResponse.
func (r BatchPublishResponse) WithError(err error) BatchPublishResponse {
	return BatchPublishResponse{
		Error:    err,
		Statuses: r.Statuses,
	}
}
