package bucketclient_test

import (
	"testing"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/scality/bucketclient/go"
)

func TestBucketclient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bucketclient Suite")
}

var client *bucketclient.BucketClient

var _ = BeforeSuite(func() {
	httpmock.Activate()
	client = bucketclient.New("http://localhost:9000")
})

var _ = AfterSuite(func() {
	httpmock.DeactivateAndReset()
})

var _ = AfterEach(func() {
	httpmock.Reset()
})
