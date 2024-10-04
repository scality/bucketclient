package bucketclient_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jarcoal/httpmock"

	"github.com/scality/bucketclient/go"
)

var testAdminGetAllSessionsInfoBucketdResponse = `[
  {
    "id": 1,
    "raftMembers": [
      {
        "id": 10,
        "name": "md1-cluster1",
        "display_name": "127.0.0.1 (md1-cluster1)",
        "host": "127.0.0.1",
        "port": 4201,
        "adminPort": 4251,
        "mdClusterId": "1"
      }
    ],
    "connectedToLeader": true
  },
  {
    "id": 2,
    "raftMembers": [
      {
        "id": 20,
        "name": "md1-cluster1",
        "display_name": "127.0.0.1 (md1-cluster1)",
        "host": "127.0.0.1",
        "port": 4202,
        "adminPort": 4252,
        "mdClusterId": "1"
      }
    ],
    "connectedToLeader": false
  }
]
`

var _ = Describe("AdminGetSessionInfo()/AdminGetAllSessionsInfo()", func() {
	Describe("AdminGetAllSessionsInfo()", func() {
		It("return info about all raft sessions", func(ctx SpecContext) {
			httpmock.RegisterResponder(
				"GET", "http://localhost:9000/_/raft_sessions",
				httpmock.NewStringResponder(200, testAdminGetAllSessionsInfoBucketdResponse),
			)
			Expect(client.AdminGetAllSessionsInfo(ctx)).To(Equal([]bucketclient.SessionInfo{
				{
					ID: 1,
					RaftMembers: []bucketclient.MemberInfo{
						{
							ID:          10,
							Name:        "md1-cluster1",
							DisplayName: "127.0.0.1 (md1-cluster1)",
							Host:        "127.0.0.1",
							Port:        4201,
							AdminPort:   4251,
							MDClusterId: "1",
						},
					},
					ConnectedToLeader: true,
				},
				{
					ID: 2,
					RaftMembers: []bucketclient.MemberInfo{
						{
							ID:          20,
							Name:        "md1-cluster1",
							DisplayName: "127.0.0.1 (md1-cluster1)",
							Host:        "127.0.0.1",
							Port:        4202,
							AdminPort:   4252,
							MDClusterId: "1",
						},
					},
					ConnectedToLeader: false,
				},
			}))
		})
	})

	Describe("AdminGetSessionInfo()", func() {
		It("return info about a particular raft session", func(ctx SpecContext) {
			httpmock.RegisterResponder(
				"GET", "http://localhost:9000/_/raft_sessions",
				httpmock.NewStringResponder(200, testAdminGetAllSessionsInfoBucketdResponse),
			)
			Expect(client.AdminGetSessionInfo(ctx, 2)).To(Equal(&bucketclient.SessionInfo{
				ID: 2,
				RaftMembers: []bucketclient.MemberInfo{
					{
						ID:          20,
						Name:        "md1-cluster1",
						DisplayName: "127.0.0.1 (md1-cluster1)",
						Host:        "127.0.0.1",
						Port:        4202,
						AdminPort:   4252,
						MDClusterId: "1",
					},
				},
				ConnectedToLeader: false,
			}))
		})
		It("return an error if the session doesn't exist", func(ctx SpecContext) {
			httpmock.RegisterResponder(
				"GET", "http://localhost:9000/_/raft_sessions",
				httpmock.NewStringResponder(200, testAdminGetAllSessionsInfoBucketdResponse),
			)
			_, err := client.AdminGetSessionInfo(ctx, 3)
			bcErr, ok := err.(*bucketclient.BucketClientError)
			Expect(ok).To(BeTrue())
			Expect(bcErr.StatusCode).To(Equal(404))
			Expect(bcErr.ErrorType).To(Equal("RaftSessionNotFound"))
		})
	})
})
