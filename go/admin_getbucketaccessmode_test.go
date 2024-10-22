package bucketclient_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jarcoal/httpmock"

	"github.com/scality/bucketclient/go"
)

var _ = Describe("AdminGetBucketAccessMode()", func() {
	It("return the access mode of the bucket", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"GET", "http://localhost:9000/_/buckets/my-bucket/accessMode",
			httpmock.NewStringResponder(200, "read-only"),
		)
		Expect(client.AdminGetBucketAccessMode(ctx, "my-bucket")).To(
			Equal(bucketclient.BucketAccessModeReadOnly))
	})
	It("return an error if the bucket doesn't exist", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"GET", "http://localhost:9000/_/buckets/nosuchbucket/accessMode",
			httpmock.NewStringResponder(404, ""),
		)
		_, err := client.AdminGetBucketAccessMode(ctx, "nosuchbucket")
		Expect(err).To(MatchError(ContainSubstring("bucketd returned HTTP status 404")))
	})
	It("escapes a bucket name containing slashes", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"GET", "http://localhost:9000/_/buckets/my-bucket%2Fwith-a-slash/accessMode",
			httpmock.NewStringResponder(200, "read-only"),
		)
		Expect(client.AdminGetBucketAccessMode(ctx, "my-bucket/with-a-slash")).To(
			Equal(bucketclient.BucketAccessModeReadOnly))
	})
})
