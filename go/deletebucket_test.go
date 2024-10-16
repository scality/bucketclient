package bucketclient_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jarcoal/httpmock"

	"github.com/scality/bucketclient/go"
)

var _ = Describe("DeleteBucket()", func() {
	It("deletes an existing bucket entry", func(ctx SpecContext) {
		httpmock.RegisterResponder("DELETE", "/default/bucket/my-bucket",
			httpmock.NewStringResponder(200, ""))

		Expect(client.DeleteBucket(ctx, "my-bucket")).To(Succeed())
	})

	It("forwards a 404 NotFound error when the bucket doesn't exist", func(ctx SpecContext) {
		httpmock.RegisterResponder("DELETE", "/default/bucket/my-bucket",
			httpmock.NewStringResponder(http.StatusNotFound, ""))

		err := client.DeleteBucket(ctx, "my-bucket")
		Expect(err).To(HaveOccurred())
		bcErr, ok := err.(*bucketclient.BucketClientError)
		Expect(ok).To(BeTrue())
		Expect(bcErr.StatusCode).To(Equal(http.StatusNotFound))
	})
})
