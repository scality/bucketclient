package bucketclient

import (
	"context"
	"fmt"
	"net/url"
)

// AdminGetBucketDBBackend returns the database backend type of the given bucket.
// Currently, it returns either "leveldb" or "rocksdb".
// Returns an empty string and an error if the bucket doesn't exist, or if a request
// error occurs.
func (client *BucketClient) AdminGetBucketDBBackend(ctx context.Context, bucketName string,
	opts ...RequestOption) (string, error) {
	// Escape the bucket name to avoid any risk to inadvertently or maliciously
	// call another route with an incorrect/crafted bucket name containing slashes.
	resource := fmt.Sprintf("/_/buckets/%s/dbBackend", url.PathEscape(bucketName))
	responseBody, err := client.Request(ctx, "AdminGetBucketDBBackend", "GET", resource, opts...)
	if err != nil {
		bcErr := err.(*BucketClientError)
		// for backward-compatibility with older Metadata versions: if the route
		// doesn't exist, assume the database type is "leveldb"
		if bcErr.StatusCode == 400 {
			return "leveldb", nil
		}
		return "", err
	}
	dbBackend := string(responseBody)
	return dbBackend, nil
}
