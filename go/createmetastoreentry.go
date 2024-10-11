package bucketclient

import (
	"context"
	"encoding/json"
	"fmt"
)

// CreateMetastoreEntry creates or updates a metastore entry for the given bucket
func (client *BucketClient) CreateMetastoreEntry(ctx context.Context, bucketName string,
	metastoreEntry MetastoreEntry) error {
	resource := fmt.Sprintf("/default/metastore/db/%s", bucketName)
	postBody, err := json.Marshal(metastoreEntry)
	if err != nil {
		return &BucketClientError{
			"CreateMetastoreEntry", "POST", client.Endpoint, resource, 0, "",
			fmt.Errorf("error marshaling POST request body: %w", err),
		}
	}
	_, err = client.Request(ctx, "CreateMetastoreEntry", "POST", resource,
		RequestBodyOption(postBody))
	return err
}
