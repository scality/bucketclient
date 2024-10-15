package bucketclient

import (
	"context"
	"fmt"
	"net/url"
)

// AdminSetBucketAccessMode sets the access mode of the given bucket:
//   - "read-write" to restore the default read/write access
//   - "read-only" to set the bucket in read-only mode and refuse write
//     operations with a 503 ServiceUnavailable error.
//
// Returns an error if the bucket doesn't exist, or if a request error occurs.
func (client *BucketClient) AdminSetBucketAccessMode(ctx context.Context,
	bucketName string, accessMode BucketAccessMode) error {
	resource := fmt.Sprintf("/_/buckets/%s/accessMode?mode=%s",
		url.PathEscape(bucketName), accessMode)
	_, err := client.Request(ctx, "AdminSetBucketAccessMode", "PUT", resource)
	return err
}
