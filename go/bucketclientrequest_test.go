package bucketclient_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jarcoal/httpmock"

	"github.com/scality/bucketclient/go"
)

type shortReadCloser struct{}

func (rc *shortReadCloser) Read(p []byte) (int, error) {
	return 0, errors.New("provoked short read")
}
func (rc *shortReadCloser) Close() error {
	return nil
}

var _ = Describe("BucketClient.Request()", func() {
	It("with invalid URL", func(ctx SpecContext) {
		// Here we do a basic test checking that errors from
		// the connection layer are forwarded correctly, "no
		// responder found" is the error returned by httpmock
		invalidClient := bucketclient.New("http://invalid:9000")
		_, err := invalidClient.Request(ctx, "GetSomething", "GET", "/foo/bar")
		Expect(err).To(MatchError(ContainSubstring("no responder found")))
	})

	Context("with valid URL", Ordered, func() {
		It("succeeds with a 200 response on GET request", func(ctx SpecContext) {
			httpmock.RegisterResponder(
				"GET", "http://localhost:9000/default/bucket/somebucket/someobject",
				httpmock.NewStringResponder(200, `{"some":"metadata","version":"1234"}`),
			)
			Expect(client.Request(ctx, "GetObject", "GET",
				"/default/bucket/somebucket/someobject")).To(Equal(
				[]byte(`{"some":"metadata","version":"1234"}`)))
		})
		It("sends PUT request with body and succeeds with a 200 response", func(ctx SpecContext) {
			httpmock.RegisterResponder(
				"PUT", "http://localhost:9000/default/bucket/somebucket/someobject",
				func(req *http.Request) (*http.Response, error) {
					defer req.Body.Close()
					Expect(io.ReadAll(req.Body)).To(Equal(
						[]byte(`{"some":"metadata","version":"1234"}`)))

					contentType, hasHeader := req.Header["Content-Type"]
					Expect(hasHeader).To(BeTrue())
					Expect(contentType).To(Equal([]string{"application/json"}))

					_, hasHeader = req.Header["Idempotency-Key"]
					Expect(hasHeader).To(BeFalse())

					return httpmock.NewStringResponse(200, "got it"), nil
				},
			)
			Expect(client.Request(ctx,
				"PutObject", "PUT", "/default/bucket/somebucket/someobject",
				bucketclient.RequestBodyContentTypeOption("application/json"),
				bucketclient.RequestBodyOption(
					[]byte(`{"some":"metadata","version":"1234"}`),
				),
			)).To(Equal([]byte("got it")))
		})
		It("sends POST request with body and idempotent flag with a 200 response", func(ctx SpecContext) {
			httpmock.RegisterResponder("POST", "http://localhost:9000/idempotent/url",
				func(req *http.Request) (*http.Response, error) {
					defer req.Body.Close()
					Expect(io.ReadAll(req.Body)).To(Equal([]byte("post body")))

					// Note: httpmock transmits the empty idempotency header, but
					// it doesn't get transmitted over the wire to bucketd
					idempotencyHeader, hasHeader := req.Header["Idempotency-Key"]
					Expect(hasHeader).To(BeTrue())
					Expect(idempotencyHeader).To(Equal([]string{}))

					return httpmock.NewStringResponse(200, "got it"), nil
				},
			)
			Expect(client.Request(ctx,
				"PostIdempotent", "POST", "/idempotent/url",
				bucketclient.RequestBodyContentTypeOption("text/plain"),
				bucketclient.RequestBodyOption([]byte("post body")),
				bucketclient.RequestIdempotent,
			)).To(Equal([]byte("got it")))
		})
		It("fails with a 400 response on GET request", func(ctx SpecContext) {
			httpmock.RegisterResponder("GET", "http://localhost:9000/invalid/url",
				httpmock.NewStringResponder(400, "dunno what to do with this"),
			)
			_, err := client.Request(ctx,
				"GetSomething", "GET", "/invalid/url")
			Expect(err).To(MatchError(ContainSubstring("GetSomething")))
			Expect(err).To(MatchError(ContainSubstring("400")))
		})
		It("fails with connection closed while reading response", func(ctx SpecContext) {
			httpmock.RegisterResponder("GET", "http://localhost:9000/some/url",
				func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						Status:        "OK",
						StatusCode:    200,
						Body:          &shortReadCloser{},
						Header:        http.Header{},
						ContentLength: 1000000,
					}, nil
				},
			)
			_, err := client.Request(ctx, "GetSomething", "GET", "/some/url")
			Expect(err).To(MatchError(ContainSubstring("GetSomething")))
			Expect(err).To(MatchError(ContainSubstring("provoked short read")))
		})
		It("fails when the context is canceled", func(ctx SpecContext) {
			httpmock.RegisterResponder(
				"GET", "http://localhost:9000/foo/bar",
				func(req *http.Request) (*http.Response, error) {
					time.Sleep(10 * time.Second)
					return httpmock.NewStringResponse(
						200, `{"some":"metadata","version":"1234"}`), nil
				},
			)
			timeoutCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
			defer cancel()

			_, err := client.Request(timeoutCtx, "GetObject", "GET", "/foo/bar")
			Expect(err).To(MatchError(context.DeadlineExceeded))
		})
	})
})
