package bucketclient

import (
	"context"
	"encoding/json"
	"fmt"
)

type Overhead struct {
	InternalOp bool `json:"internalOp"`
}

type PostBatchEntry struct {
	Key      string    `json:"key"`
	Value    string    `json:"value,omitempty"`
	Overhead *Overhead `json:"overhead,omitempty"`
	Type     string    `json:"type,omitempty"`
}

func (client *BucketClient) PostBatch(ctx context.Context,
	bucketName string, batch []PostBatchEntry, opts ...RequestOption) error {
	resource := fmt.Sprintf("/default/batch/%s", bucketName)
	postPayload := struct {
		Batch []PostBatchEntry `json:"batch"`
	}{Batch: batch}
	postBody, err := json.Marshal(postPayload)
	fmt.Println("marshal", postBody, err)
	if err != nil {
		return &BucketClientError{
			"PostBatch", "POST", client.Endpoint, resource, 0, "",
			fmt.Errorf("error marshaling POST request body: %w", err),
		}
	}
	_, err = client.Request(ctx, "PostBatch", "POST", resource,
		append([]RequestOption{
			RequestBodyOption(postBody),
			RequestBodyContentTypeOption("application/json"),
			// Because we write a batch of low-level entries directly to
			// the database, the request is idempotent.
			RequestIdempotent,
		}, opts...)...)
	fmt.Println("postbatch", err)
	return err
}
