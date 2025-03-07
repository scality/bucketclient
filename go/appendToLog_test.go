package bucketclient_test

import (
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"io"
	"net/http"

	bucketclient "github.com/scality/bucketclient/go"
)

var _ = Describe("AppendToLog()", func() {
	It("POSTs a batch to bucketd's appendToLog route", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"POST", "http://localhost:9000/default/appendToLog/somebucket",
			func(req *http.Request) (*http.Response, error) {
				defer req.Body.Close()
				Expect(io.ReadAll(req.Body)).To(Equal(
					[]byte(`{"batch":[{"type":"append","value":"{}"}]}`)))

				contentType, hasHeader := req.Header["Content-Type"]
				Expect(hasHeader).To(BeTrue())
				Expect(contentType).To(Equal([]string{"application/json"}))

				return httpmock.NewStringResponse(200, "got it"), nil
			},
		)

		Expect(client.AppendToLog(ctx, "somebucket", []bucketclient.AppendToLogEntry{
			{Type: "append", Value: "{}"},
		})).To(Succeed())
	})
})
