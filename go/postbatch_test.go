package bucketclient_test

import (
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"io"
	"net/http"

	"github.com/scality/bucketclient/go"
)

var _ = Describe("PostBatch()", func() {
	It("POSTs a batch to bucketd", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"POST", "http://localhost:9000/default/batch/somebucket",
			func(req *http.Request) (*http.Response, error) {
				defer req.Body.Close()
				Expect(io.ReadAll(req.Body)).To(Equal(
					[]byte(`{"batch":[{"key":"foo","value":"{}"},{"key":"bar","type":"del"}]}`)))

				contentType, hasHeader := req.Header["Content-Type"]
				Expect(hasHeader).To(BeTrue())
				Expect(contentType).To(Equal([]string{"application/json"}))

				return httpmock.NewStringResponse(200, "got it"), nil
			},
		)

		Expect(client.PostBatch(ctx, "somebucket", []bucketclient.PostBatchEntry{
			{Key: "foo", Value: "{}"},
			{Key: "bar", Type: "del"},
		})).To(Succeed())
	})
})
