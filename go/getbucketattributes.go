package bucketclient

import (
	"context"
	"fmt"
)

// GetBucketAttributes retrieves the JSON blob containing the bucket
// attributes attached to a bucket.
func (client *BucketClient) GetBucketAttributes(ctx context.Context, bucketName string) ([]byte, error) {
	resource := fmt.Sprintf("/default/attributes/%s", bucketName)
	return client.Request(ctx, "GetBucketAttributes", "GET", resource)
}
