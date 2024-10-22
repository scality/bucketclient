package bucketclient_test

import (
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"io"
	"net/http"

	"github.com/scality/bucketclient/go"
)

var _ = Describe("CreateMetastoreEntry()", func() {
	It("creates/updates a new metastore entry", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"POST", "/default/metastore/db/my-bucket",
			func(req *http.Request) (*http.Response, error) {
				defer req.Body.Close()
				Expect(io.ReadAll(req.Body)).To(Equal(
					[]byte(`{"name":"bucket-03","attributes":"{\"name\":\"bucket-03\"}",` +
						`"creating":false,"deleting":false,"id":"bucket-03",` +
						`"raftSessionID":1,"version":2,"raftSession":"rs-1"}`)))
				return httpmock.NewStringResponse(200, ""), nil
			},
		)
		Expect(client.CreateMetastoreEntry(ctx,
			"my-bucket", bucketclient.MetastoreEntry{
				Name:          "bucket-03",
				Attributes:    `{"name":"bucket-03"}`,
				Creating:      false,
				Deleting:      false,
				ID:            "bucket-03",
				RaftSessionID: 1,
				Version:       2,
				RaftSession:   "rs-1",
			})).To(Succeed())
	})
})
