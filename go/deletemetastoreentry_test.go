package bucketclient_test

import (
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/scality/bucketclient/go"
)

var _ = Describe("DeleteMetastoreEntry()", func() {
	It("deletes a metastore entry", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"DELETE", "/default/metastore/db/my-bucket",
			httpmock.NewStringResponder(200, ""),
		)
		Expect(client.DeleteMetastoreEntry(ctx, "my-bucket")).To(Succeed())
	})

	It("return 404 error if bucket does not exist", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"DELETE", "/default/metastore/db/doesnotexist",
			httpmock.NewStringResponder(404, ""),
		)
		err := client.DeleteMetastoreEntry(ctx, "doesnotexist")
		Expect(err).To(HaveOccurred())
		bcErr, ok := err.(*bucketclient.BucketClientError)
		Expect(ok).To(BeTrue())
		Expect(bcErr.StatusCode).To(Equal(404))
	})
})
