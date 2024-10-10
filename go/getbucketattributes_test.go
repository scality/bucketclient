package bucketclient_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jarcoal/httpmock"

	"github.com/scality/bucketclient/go"
)

var _ = Describe("GetBucketAttributes()", func() {
	It("retrieves the bucket attributes of an existing bucket", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"GET", "/default/attributes/my-bucket",
			httpmock.NewStringResponder(200, `{"foo":"bar"}`),
		)
		Expect(client.GetBucketAttributes(ctx, "my-bucket")).To(
			Equal([]byte(`{"foo":"bar"}`)))
	})

	It("returns a 404 error if the bucket does not exist", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"GET", "/default/attributes/my-bucket",
			httpmock.NewStringResponder(404, ""),
		)
		_, err := client.GetBucketAttributes(ctx, "my-bucket")
		Expect(err).To(HaveOccurred())
		bcErr, ok := err.(*bucketclient.BucketClientError)
		Expect(ok).To(BeTrue())
		Expect(bcErr.StatusCode).To(Equal(404))
	})
})
