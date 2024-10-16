package bucketclient

import (
	"context"
	"fmt"
)

// DeleteBucket deletes a bucket entry from metadata.
func (client *BucketClient) DeleteBucket(ctx context.Context, bucketName string) error {
	resource := fmt.Sprintf("/default/bucket/%s", bucketName)

	_, err := client.Request(ctx, "DeleteBucket", "DELETE", resource)
	return err
}
