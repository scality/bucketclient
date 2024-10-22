package bucketclient

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// AdminGetBucketSessionID returns the raft session ID of the given bucket.
// Returns 0 and an error if the bucket doesn't exist, or if a request error occurs.
func (client *BucketClient) AdminGetBucketSessionID(ctx context.Context, bucketName string) (int, error) {
	// Escape the bucket name to avoid any risk to inadvertently or maliciously
	// call another route with an incorrect/crafted bucket name containing slashes.
	resource := fmt.Sprintf("/_/buckets/%s/id", url.PathEscape(bucketName))
	responseBody, err := client.Request(ctx, "AdminGetBucketSessionID", "GET", resource)
	if err != nil {
		return 0, err
	}
	sessionId, err := strconv.ParseInt(string(responseBody), 10, 0)
	if err != nil {
		return 0, &BucketClientError{
			"AdminGetBucketSessionID", "GET", client.Endpoint, resource, 0, "",
			fmt.Errorf("bucketd did not return a valid session ID in response body: '%s'",
				string(responseBody)),
		}
	}
	return int(sessionId), nil
}
