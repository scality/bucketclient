package bucketclient_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jarcoal/httpmock"

	"net/http"

	"github.com/scality/bucketclient/go"
)

var _ = Describe("AdminGetSessionLog()", func() {
	It("returns a range of raft oplog", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"GET", "http://localhost:9000/_/raft_sessions/2/log?begin=10&limit=2",
			httpmock.NewStringResponder(200, `{
  "info": {
    "start": 10,
    "cseq": 100,
    "prune": 1
  },
  "log": [
    {
      "db": "bucket1",
      "method": 8,
      "entries": [
        {
          "key": "key1",
          "value": "{\"foo\":\"bar\"}"
        },
        {
          "key": "key2",
          "type": "del"
        }
      ]
    },
    {
      "db": "bucket2",
      "method": 0,
      "entries": [
        {
          "value": "{}"
        }
      ]
    }
  ]
}
`),
		)
		Expect(
			client.AdminGetSessionLog(ctx, 2, 10, 2, false),
		).To(Equal(&bucketclient.AdminGetSessionLogResponse{
			Info: bucketclient.SessionLogInfo{
				Start: 10,
				CSeq:  100,
				Prune: 1,
			},
			Log: []bucketclient.SessionLogRecord{
				bucketclient.SessionLogRecord{
					Bucket:    "bucket1",
					DBMethod:  bucketclient.DBMethodBatch,
					Timestamp: "",
					Entries: []bucketclient.SessionLogEntry{
						bucketclient.SessionLogEntry{
							Key:   "key1",
							Value: "{\"foo\":\"bar\"}",
							Type:  "",
						},
						bucketclient.SessionLogEntry{
							Key:   "key2",
							Value: "",
							Type:  "del",
						},
					},
				},
				bucketclient.SessionLogRecord{
					Bucket:    "bucket2",
					DBMethod:  bucketclient.DBMethodCreate,
					Timestamp: "",
					Entries: []bucketclient.SessionLogEntry{
						bucketclient.SessionLogEntry{Key: "", Value: "{}", Type: ""},
					},
				},
			},
		}))
	})

	It("returns an error with status 416 RequestedRangeNotSatisfiable", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"GET", "http://localhost:9000/_/raft_sessions/2/log?begin=10&limit=2",
			httpmock.NewStringResponder(http.StatusRequestedRangeNotSatisfiable, ""),
		)
		_, err := client.AdminGetSessionLog(ctx, 2, 10, 2, false)
		Expect(err).To(HaveOccurred())
		bcErr, ok := err.(*bucketclient.BucketClientError)
		Expect(ok).To(BeTrue())
		Expect(bcErr.StatusCode).To(Equal(http.StatusRequestedRangeNotSatisfiable))
	})
})
