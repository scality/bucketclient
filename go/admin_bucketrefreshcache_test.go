package bucketclient_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jarcoal/httpmock"

	"github.com/scality/bucketclient/go"
)

var _ = Describe("AdminBucketRefreshCache()", func() {
	It("calls the refreshCache route", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"GET", "http://localhost:9000/_/buckets/my-bucket/refreshCache",
			httpmock.NewStringResponder(200, ""),
		)
		Expect(client.AdminBucketRefreshCache(ctx, "my-bucket")).To(Succeed())
	})
	It("return an error if the bucket doesn't exist", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"GET", "http://localhost:9000/_/buckets/nosuchbucket/refreshCache",
			httpmock.NewStringResponder(404, ""),
		)
		err := client.AdminBucketRefreshCache(ctx, "nosuchbucket")
		Expect(err).To(HaveOccurred())
		bcErr, ok := err.(*bucketclient.BucketClientError)
		Expect(ok).To(BeTrue())
		Expect(bcErr.StatusCode).To(Equal(404))
	})
	It("escapes a bucket name containing slashes", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"GET", "http://localhost:9000/_/buckets/my-bucket%2Fwith-a-slash/refreshCache",
			httpmock.NewStringResponder(200, ""),
		)
		Expect(client.AdminBucketRefreshCache(ctx, "my-bucket/with-a-slash")).To(Succeed())
	})
})
