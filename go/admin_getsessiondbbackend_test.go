package bucketclient_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jarcoal/httpmock"
)

var _ = Describe("AdminGetSessionDBBackend()", func() {
	It("return the database backend type of the requested RAFT session", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"GET", "http://localhost:9000/_/raft_sessions/4/dbBackend",
			httpmock.NewStringResponder(200, "rocksdb"),
		)
		Expect(client.AdminGetSessionDBBackend(ctx, 4)).To(Equal("rocksdb"))
	})
	It("forwards request errors", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"GET", "http://localhost:9000/_/raft_sessions/4/dbBackend",
			httpmock.NewStringResponder(500, ""),
		)
		_, err := client.AdminGetSessionDBBackend(ctx, 4)
		Expect(err).To(MatchError(ContainSubstring("bucketd returned HTTP status 500")))
	})
	It("return 'leveldb' if the route doesn't exist (backward-compatibility)", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"GET", "http://localhost:9000/_/raft_sessions/4/dbBackend",
			httpmock.NewStringResponder(400, "BadRequest"),
		)
		Expect(client.AdminGetSessionDBBackend(ctx, 4)).To(Equal("leveldb"))
	})
})
