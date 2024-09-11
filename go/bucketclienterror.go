package bucketclient

import (
	"fmt"
)

type BucketClientError struct {
	ApiMethod  string
	HttpMethod string
	Endpoint   string
	Resource   string
	StatusCode int
	ErrorType  string
	Err        error
}

func (e *BucketClientError) Error() string {
	if e.StatusCode > 0 {
		return fmt.Sprintf("error in %s [%s %s%s]: bucketd returned HTTP status %d %s",
			e.ApiMethod, e.HttpMethod, e.Endpoint, e.Resource, e.StatusCode, e.ErrorType)
	} else {
		return fmt.Sprintf("error in %s [%s %s%s]: HTTP request to bucketd failed: %v",
			e.ApiMethod, e.HttpMethod, e.Endpoint, e.Resource, e.Err)
	}
}

func (e *BucketClientError) Unwrap() error {
	return e.Err
}

func ErrorMalformedResponse(apiMethod string, httpMethod string, endpoint string, resource string,
	err error) error {
	return &BucketClientError{apiMethod, httpMethod, endpoint, resource, 0, "",
		fmt.Errorf("bucketd returned a malformed response body: %w", err),
	}
}
