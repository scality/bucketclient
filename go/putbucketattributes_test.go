package bucketclient_test

import (
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"io"
	"net/http"
)

var _ = Describe("PutBucketAttributes()", func() {
	It("updates bucket attributes with a new JSON blob", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"POST", "/default/attributes/my-bucket",
			func(req *http.Request) (*http.Response, error) {
				defer req.Body.Close()
				Expect(io.ReadAll(req.Body)).To(Equal([]byte(`{"foo":"bar"}`)))
				return httpmock.NewStringResponse(200, ""), nil
			},
		)
		Expect(client.PutBucketAttributes(ctx,
			"my-bucket", []byte(`{"foo":"bar"}`))).To(Succeed())
	})
})
