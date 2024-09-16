package bucketclient_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jarcoal/httpmock"

	"github.com/scality/bucketclient/go"
)

var testAdminGetSessionLeaderBucketdResponse = `{
  "id": 10,
  "name": "md1-cluster1",
  "display_name": "127.0.0.1 (md1-cluster1)",
  "host": "127.0.0.1",
  "port": 4201,
  "adminPort": 4251,
  "mdClusterId": "1"
}
`

var _ = Describe("AdminGetSessionLeader()", func() {
	It("returns a MemberInfo about the leader of a raft session", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"GET", "http://localhost:9000/_/raft_sessions/3/leader",
			httpmock.NewStringResponder(200, testAdminGetSessionLeaderBucketdResponse),
		)
		Expect(client.AdminGetSessionLeader(ctx, 3)).To(Equal(&bucketclient.MemberInfo{
			ID:          10,
			Name:        "md1-cluster1",
			DisplayName: "127.0.0.1 (md1-cluster1)",
			Host:        "127.0.0.1",
			Port:        4201,
			AdminPort:   4251,
			MDClusterId: "1",
		}))
	})

	It("forwards an error from bucketd", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"GET", "http://localhost:9000/_/raft_sessions/3/leader",
			httpmock.NewStringResponder(400, ""),
		)
		_, err := client.AdminGetSessionLeader(ctx, 3)
		Expect(err).To(HaveOccurred())
	})
})
