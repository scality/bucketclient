package bucketclient

import (
	"net/http"
)

type BucketClient struct {
	*http.Client
	Endpoint string
	Metrics  *BucketClientMetrics
}

// New creates a new BucketClient instance, with the provided endpoint (e.g. "localhost:9000")
// and the default HTTP client.
func New(bucketdEndpoint string) *BucketClient {
	bc := &BucketClient{
		Client:   http.DefaultClient,
		Endpoint: bucketdEndpoint,
	}
	return bc
}

// NewWithHTTPClient creates a new BucketClient instance, with the provided endpoint
// (e.g. "localhost:9000") using the provided http.Client instance.
func NewWithHTTPClient(bucketdEndpoint string, httpClient *http.Client) *BucketClient {
	return &BucketClient{
		Client:   httpClient,
		Endpoint: bucketdEndpoint,
	}
}
