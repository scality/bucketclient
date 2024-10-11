package bucketclient_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jarcoal/httpmock"

	"github.com/scality/bucketclient/go"
)

var testGetMetastoreEntryResponse = `{"name":"bucket-03","attributes":"{\"name\":\"bucket-03\"}",` +
	`"creating":false,"deleting":false,"id":"bucket-03","raftSessionID":1,"version":2,` +
	`"raftSession":"rs-1"}`

var _ = Describe("GetMetastoreEntry()", func() {
	It("retrieves the parsed metastore entry for a bucket", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"GET", "/default/metastore/db/my-bucket",
			httpmock.NewStringResponder(200, testGetMetastoreEntryResponse),
		)
		Expect(client.GetMetastoreEntry(ctx, "my-bucket")).To(
			Equal(bucketclient.MetastoreEntry{
				Name:          "bucket-03",
				Attributes:    `{"name":"bucket-03"}`,
				Creating:      false,
				Deleting:      false,
				ID:            "bucket-03",
				RaftSessionID: 1,
				Version:       2,
				RaftSession:   "rs-1",
			}))
	})

	It("returns a 404 error if the bucket does not exist", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"GET", "/default/metastore/db/my-bucket",
			httpmock.NewStringResponder(404, ""),
		)
		_, err := client.GetMetastoreEntry(ctx, "my-bucket")
		Expect(err).To(HaveOccurred())
		bcErr, ok := err.(*bucketclient.BucketClientError)
		Expect(ok).To(BeTrue())
		Expect(bcErr.StatusCode).To(Equal(404))
	})
})
