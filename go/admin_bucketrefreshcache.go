package bucketclient

import (
	"context"
	"fmt"
	"net/url"
)

// AdminBucketRefreshCache refreshes the bucketd cache of metastore
// entries for the given bucket. Useful after switching the raft session
// of a bucket.
func (client *BucketClient) AdminBucketRefreshCache(ctx context.Context, bucketName string) error {
	resource := fmt.Sprintf("/_/buckets/%s/refreshCache", url.PathEscape(bucketName))
	_, err := client.Request(ctx, "AdminBucketRefreshCache", "GET", resource)
	return err
}
