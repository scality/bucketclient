package bucketclient

type BucketClient struct {
	Endpoint string
}

func New(bucketdEndpoint string) *BucketClient {
	return &BucketClient{bucketdEndpoint}
}
