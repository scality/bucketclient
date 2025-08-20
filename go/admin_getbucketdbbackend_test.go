package bucketclient_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jarcoal/httpmock"
)

var _ = Describe("AdminGetBucketDBBackend()", func() {
	It("return the database backend type of the RAFT session hosting a bucket", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"GET", "http://localhost:9000/_/buckets/my-bucket/dbBackend",
			httpmock.NewStringResponder(200, "rocksdb"),
		)
		Expect(client.AdminGetBucketDBBackend(ctx, "my-bucket")).To(Equal("rocksdb"))
	})
	It("return an error if the bucket doesn't exist", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"GET", "http://localhost:9000/_/buckets/nosuchbucket/dbBackend",
			httpmock.NewStringResponder(404, ""),
		)
		_, err := client.AdminGetBucketDBBackend(ctx, "nosuchbucket")
		Expect(err).To(MatchError(ContainSubstring("bucketd returned HTTP status 404")))
	})
	It("return 'leveldb' if the route doesn't exist (backward-compatibility)", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"GET", "http://localhost:9000/_/buckets/my-bucket/dbBackend",
			httpmock.NewStringResponder(400, "BadRequest"),
		)
		Expect(client.AdminGetBucketDBBackend(ctx, "my-bucket")).To(Equal("leveldb"))
	})
	It("escapes a bucket name containing slashes", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"GET", "http://localhost:9000/_/buckets/my-bucket%2Fwith-a-slash/dbBackend",
			httpmock.NewStringResponder(200, "rocksdb"),
		)
		Expect(client.AdminGetBucketDBBackend(ctx, "my-bucket/with-a-slash")).To(Equal("rocksdb"))
	})
})
