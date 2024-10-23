package bucketclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type createBucketOptionSet struct {
	sessionId      int
	makeIdempotent bool
}

type CreateBucketOption func(*createBucketOptionSet)

func CreateBucketSessionIdOption(sessionId int) CreateBucketOption {
	return func(options *createBucketOptionSet) {
		options.sessionId = sessionId
	}
}

func CreateBucketMakeIdempotent(options *createBucketOptionSet) {
	options.makeIdempotent = true
}

// CreateBucket creates a bucket in metadata.
// bucketAttributes is a JSON blob of bucket attributes
// opts is a set of options:
//
//	CreateBucketSessionIdOption forces the session ID where the bucket to be
//	    created will land
//
//	CreateBucketMakeIdempotent makes the request return a success if a bucket
//	    with the same UID already exists (otherwise returns 409 Conflict, as
//	    if the option is not passed)
func (client *BucketClient) CreateBucket(ctx context.Context,
	bucketName string, bucketAttributes []byte, opts ...CreateBucketOption) error {
	parsedOpts := createBucketOptionSet{
		sessionId:      0,
		makeIdempotent: false,
	}
	for _, opt := range opts {
		opt(&parsedOpts)
	}

	resource := fmt.Sprintf("/default/bucket/%s", bucketName)
	query := url.Values{}

	if parsedOpts.sessionId > 0 {
		query.Set("raftsession", strconv.Itoa(parsedOpts.sessionId))
	}
	u, _ := url.Parse(resource)
	u.RawQuery = query.Encode()
	resource = u.String()
	requestOptions := []RequestOption{
		RequestBodyOption(bucketAttributes),
		RequestBodyContentTypeOption("application/json"),
	}
	if parsedOpts.makeIdempotent {
		// since we will make the request idempotent, it's
		// okay to retry it (it may return 409 Conflict at the
		// first retry if it initially succeeded, but it will
		// then be considered a success)
		requestOptions = append(requestOptions, RequestIdempotent)
	}
	_, err := client.Request(ctx, "CreateBucket", "POST", resource, requestOptions...)
	if err == nil {
		return nil
	}
	if parsedOpts.makeIdempotent {
		// If the Idempotent option is set, Accept "409 Conflict" as a success iff
		// the UIDs match between the existing and the new metadata, to detect and
		// return an error if there is an existing bucket that was not created by us

		bcErr := err.(*BucketClientError)
		if bcErr.StatusCode != http.StatusConflict {
			return err
		}
		existingBucketAttributes, err := client.GetBucketAttributes(ctx, bucketName)
		if err != nil {
			return err
		}
		if bucketAttributeUIDsMatch(bucketAttributes, existingBucketAttributes) {
			// return silent success without updating the existing metadata
			return nil
		}
	}
	return err
}

func bucketAttributeUIDsMatch(attributes1 []byte, attributes2 []byte) bool {
	var parsedAttr1, parsedAttr2 struct {
		UID string `json:"uid"`
	}

	err := json.Unmarshal(attributes1, &parsedAttr1)
	if err != nil {
		return false
	}
	err = json.Unmarshal(attributes2, &parsedAttr2)
	if err != nil {
		return false
	}
	if parsedAttr1.UID == "" || parsedAttr2.UID == "" {
		return false
	}
	return parsedAttr1.UID == parsedAttr2.UID
}
