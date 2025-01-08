package bucketclient

import (
	"context"
	"fmt"
)

// DeleteMetastoreEntry deletes the metastore entry for the given bucket
func (client *BucketClient) DeleteMetastoreEntry(ctx context.Context, bucketName string,
	opts ...RequestOption) error {
	resource := fmt.Sprintf("/default/metastore/db/%s", bucketName)
	_, err := client.Request(ctx, "DeleteMetastoreEntry", "DELETE", resource, opts...)
	return err
}
