package bucketclient_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jarcoal/httpmock"
)

var _ = Describe("AdminGetBucketSessionId()", func() {
	It("return the raft session ID hosting a bucket", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"GET", "http://localhost:9000/_/buckets/my-bucket/id",
			httpmock.NewStringResponder(200, "42"),
		)
		Expect(client.AdminGetBucketSessionID(ctx, "my-bucket")).To(Equal(42))
	})
	It("return an error if the bucket doesn't exist", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"GET", "http://localhost:9000/_/buckets/nosuchbucket/id",
			httpmock.NewStringResponder(404, ""),
		)
		_, err := client.AdminGetBucketSessionID(ctx, "nosuchbucket")
		Expect(err).To(MatchError(ContainSubstring("bucketd returned HTTP status 404")))
	})
	It("escapes a bucket name containing slashes", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"GET", "http://localhost:9000/_/buckets/my-bucket%2Fwith-a-slash/id",
			httpmock.NewStringResponder(200, "42"),
		)
		Expect(client.AdminGetBucketSessionID(ctx, "my-bucket/with-a-slash")).To(Equal(42))
	})
})
