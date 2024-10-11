package bucketclient

import (
	"context"
	"encoding/json"
	"fmt"
)

type PostBatchEntry struct {
	Key   string `json:"key"`
	Value string `json:"value,omitempty"`
	Type  string `json:"type,omitempty"`
}

func (client *BucketClient) PostBatch(ctx context.Context,
	bucketName string, batch []PostBatchEntry) error {
	resource := fmt.Sprintf("/default/batch/%s", bucketName)
	postPayload := struct {
		Batch []PostBatchEntry `json:"batch"`
	}{Batch: batch}
	postBody, err := json.Marshal(postPayload)
	if err != nil {
		return &BucketClientError{
			"PostBatch", "POST", client.Endpoint, resource, 0, "",
			fmt.Errorf("error marshaling POST request body: %w", err),
		}
	}
	_, err = client.Request(ctx, "PostBatch", "POST", resource,
		RequestBodyOption(postBody),
		RequestBodyContentTypeOption("application/json"),
		// Because we write a batch of low-level entries directly to
		// the database, the request is idempotent.
		RequestIdempotent)
	return err
}
