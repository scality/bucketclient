package bucketclient

import (
	"context"
	"fmt"
)

// AdminGetBucketAccessMode returns the access mode of the given bucket:
// - "read-write" when the bucket is accessible for reading and writing (default)
// - "read-only" when the bucket is only accessible for read operations.
// Returns "" and an error if the bucket doesn't exist, or if a request error occurs.
func (client *BucketClient) AdminGetBucketAccessMode(ctx context.Context,
	bucketName string) (BucketAccessMode, error) {
	resource := fmt.Sprintf("/_/buckets/%s/accessMode", bucketName)
	responseBody, err := client.Request(ctx, "AdminGetBucketAccessMode", "GET", resource)
	if err != nil {
		return "", err
	}
	accessMode := BucketAccessMode(string(responseBody))
	return accessMode, nil
}
