package bucketclient

import (
	"context"
	"fmt"
	"net/url"
)

// AdminGetBucketAccessMode returns the access mode of the given bucket:
// - "read-write" when the bucket is accessible for reading and writing (default)
// - "read-only" when the bucket is only accessible for read operations.
// Returns "" and an error if the bucket doesn't exist, or if a request error occurs.
func (client *BucketClient) AdminGetBucketAccessMode(ctx context.Context,
	bucketName string) (BucketAccessMode, error) {
	// Escape the bucket name to avoid any risk to inadvertently or maliciously
	// call another route with an incorrect/crafted bucket name containing slashes.
	resource := fmt.Sprintf("/_/buckets/%s/accessMode", url.PathEscape(bucketName))
	responseBody, err := client.Request(ctx, "AdminGetBucketAccessMode", "GET", resource)
	if err != nil {
		return "", err
	}
	accessMode := BucketAccessMode(string(responseBody))
	return accessMode, nil
}
