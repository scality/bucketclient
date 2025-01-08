package bucketclient

import (
	"context"
	"encoding/json"
	"fmt"
)

// GetMetastoreEntry retrieves and parses a metastore entry for the given bucket
func (client *BucketClient) GetMetastoreEntry(ctx context.Context, bucketName string,
	opts ...RequestOption) (MetastoreEntry, error) {
	resource := fmt.Sprintf("/default/metastore/db/%s", bucketName)
	responseBody, err := client.Request(ctx, "GetMetastoreEntry", "GET", resource, opts...)
	if err != nil {
		return MetastoreEntry{}, err
	}
	var metastoreEntry MetastoreEntry
	jsonErr := json.Unmarshal(responseBody, &metastoreEntry)
	if jsonErr != nil {
		return MetastoreEntry{}, ErrorMalformedResponse("GetMetastoreEntry",
			"GET", client.Endpoint, resource, jsonErr)
	}
	return metastoreEntry, nil
}
