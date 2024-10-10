package bucketclient

import (
	"context"
	"fmt"
	"net/url"
)

// PutBucketAttributes updates the bucket attributes with a new JSON blob.
func (client *BucketClient) PutBucketAttributes(ctx context.Context, bucketName string,
	bucketAttributes []byte) error {
	resource := fmt.Sprintf("/default/attributes/%s", url.PathEscape(bucketName))
	_, err := client.Request(ctx, "PutBucketAttributes", "POST", resource,
		RequestBodyOption(bucketAttributes))
	return err
}
