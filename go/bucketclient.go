package bucketclient

import (
	"net/http"
)

type BucketClient struct {
	Endpoint   string
	HTTPClient *http.Client
}

// New creates a new BucketClient instance, with the provided endpoint (e.g. "localhost:9000")
// and the default HTTP client.
func New(bucketdEndpoint string) *BucketClient {
	return &BucketClient{
		Endpoint:   bucketdEndpoint,
		HTTPClient: http.DefaultClient,
	}
}

// NewWithHTTPClient creates a new BucketClient instance, with the provided endpoint
// (e.g. "localhost:9000") using the provided http.Client instance.
func NewWithHTTPClient(bucketdEndpoint string, httpClient *http.Client) *BucketClient {
	return &BucketClient{
		Endpoint:   bucketdEndpoint,
		HTTPClient: httpClient,
	}
}
