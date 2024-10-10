package bucketclient_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jarcoal/httpmock"

	"github.com/scality/bucketclient/go"
)

var _ = Describe("AdminSetBucketAccessMode()", func() {
	It("sets the access mode of the bucket", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"PUT", "http://localhost:9000/_/buckets/my-bucket/accessMode?mode=read-only",
			httpmock.NewStringResponder(200, ""),
		)
		Expect(client.AdminSetBucketAccessMode(ctx,
			"my-bucket", bucketclient.BucketAccessModeReadOnly)).To(Succeed())
	})
	It("return an error if the bucket doesn't exist", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"PUT", "http://localhost:9000/_/buckets/nosuchbucket/accessMode?mode=read-only",
			httpmock.NewStringResponder(404, ""),
		)
		err := client.AdminSetBucketAccessMode(ctx,
			"nosuchbucket", bucketclient.BucketAccessModeReadOnly)
		Expect(err).To(MatchError(ContainSubstring("bucketd returned HTTP status 404")))
	})
	It("escapes a bucket name containing slashes", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"PUT", "http://localhost:9000/_/buckets/my-bucket%2Fwith-a-slash/accessMode?mode=read-only",
			httpmock.NewStringResponder(200, ""),
		)
		Expect(client.AdminSetBucketAccessMode(ctx,
			"my-bucket/with-a-slash", bucketclient.BucketAccessModeReadOnly)).To(Succeed())
	})
})
