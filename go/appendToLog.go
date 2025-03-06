package bucketclient

import (
	"context"
	"encoding/json"
	"fmt"
)

type AppendToLogEntry struct {
	Type     string    `json:"type"`
	Key      string    `json:"key,omitempty"`
	Value    string    `json:"value,omitempty"`
	Overhead *Overhead `json:"overhead,omitempty"`
}

func (client *BucketClient) AppendToLog(ctx context.Context,
	bucketName string, batch []AppendToLogEntry, opts ...RequestOption) error {
	resource := fmt.Sprintf("/default/appendToLog/%s", bucketName)
	postPayload := struct {
		Batch []AppendToLogEntry `json:"batch"`
	}{Batch: batch}
	postBody, err := json.Marshal(postPayload)
	if err != nil {
		return &BucketClientError{
			"AppendToLog", "POST", client.Endpoint, resource, 0, "",
			fmt.Errorf("error marshaling POST request body: %w", err),
		}
	}
	_, err = client.Request(ctx, "AppendToLog", "POST", resource,
		append([]RequestOption{
			RequestBodyOption(postBody),
			RequestBodyContentTypeOption("application/json"),
			RequestIdempotent,
		}, opts...)...)
	return err
}
