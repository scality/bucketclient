package bucketclient_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/scality/bucketclient/go"
)

var _ = Describe("BucketClient", func() {
	It("New", func() {
		client := bucketclient.New("http://localhost:9000")
		Expect(client).ToNot(BeNil())
		Expect(client.Endpoint).To(Equal("http://localhost:9000"))
	})
})
