package bucketclient

type BucketClient struct {
	Endpoint string
	Metrics  *BucketClientMetrics
}

func New(bucketdEndpoint string) *BucketClient {
	return &BucketClient{
		Endpoint: bucketdEndpoint,
	}
}
