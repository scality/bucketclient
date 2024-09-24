package bucketclient_test

import (
	"io"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jarcoal/httpmock"

	"github.com/scality/bucketclient/go"
)

var _ = Describe("CreateBucket()", func() {
	It("creates a bucket on an available raft session", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"POST", "/default/bucket/my-new-bucket",
			func(req *http.Request) (*http.Response, error) {
				defer req.Body.Close()
				Expect(io.ReadAll(req.Body)).To(Equal([]byte(`{"foo":"bar"}`)))
				return httpmock.NewStringResponse(200, ""), nil
			},
		)
		Expect(client.CreateBucket(ctx, "my-new-bucket", []byte(`{"foo":"bar"}`))).To(Succeed())
	})

	It("creates a bucket on a chosen raft session", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"POST", "/default/bucket/my-new-bucket?raftsession=12",
			func(req *http.Request) (*http.Response, error) {
				defer req.Body.Close()
				Expect(io.ReadAll(req.Body)).To(Equal([]byte(`{"foo":"bar"}`)))
				return httpmock.NewStringResponse(200, ""), nil
			},
		)
		Expect(client.CreateBucket(ctx, "my-new-bucket", []byte(`{"foo":"bar"}`),
			bucketclient.CreateBucketSessionIdOption(12))).To(Succeed())
	})

	It("forwards request error", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"POST", "/default/bucket/my-new-bucket",
			httpmock.NewStringResponder(http.StatusInternalServerError, "I'm afraid I can't do this"),
		)
		err := client.CreateBucket(ctx, "my-new-bucket", []byte(`{"foo":"bar"}`),
			bucketclient.CreateBucketMakeIdempotent)
		Expect(err).To(HaveOccurred())
		bcErr, ok := err.(*bucketclient.BucketClientError)
		Expect(ok).To(BeTrue())
		Expect(bcErr.StatusCode).To(Equal(http.StatusInternalServerError))
	})

	It("returns 409 Conflict without MakeIdempotent option if bucket with same UID exists", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"POST", "/default/bucket/my-new-bucket",
			httpmock.NewStringResponder(http.StatusConflict, ""),
		)
		// normally unused, but set to match the following tests
		httpmock.RegisterResponder(
			"GET", "/default/attributes/my-new-bucket",
			httpmock.NewStringResponder(200, `{"foo":"bar","uid":"4242"}`),
		)
		err := client.CreateBucket(ctx, "my-new-bucket", []byte(`{"foo":"bar","uid":"4242"}`))
		Expect(err).To(HaveOccurred())
		bcErr, ok := err.(*bucketclient.BucketClientError)
		Expect(ok).To(BeTrue())
		Expect(bcErr.StatusCode).To(Equal(http.StatusConflict))
	})

	It("succeeds to create bucket with MakeIdempotent option if bucket with same UID exists", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"POST", "/default/bucket/my-new-bucket",
			httpmock.NewStringResponder(http.StatusConflict, ""),
		)
		httpmock.RegisterResponder(
			"GET", "/default/attributes/my-new-bucket",
			httpmock.NewStringResponder(200, `{"foo":"bar","uid":"4242"}`),
		)
		Expect(client.CreateBucket(ctx, "my-new-bucket", []byte(`{"foo":"bar","uid":"4242"}`),
			bucketclient.CreateBucketMakeIdempotent)).To(Succeed())
	})

	It("returns 409 Conflict with MakeIdempotent option if bucket with different UID exists", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"POST", "/default/bucket/my-new-bucket",
			httpmock.NewStringResponder(http.StatusConflict, ""),
		)
		httpmock.RegisterResponder(
			"GET", "/default/attributes/my-new-bucket",
			httpmock.NewStringResponder(200, `{"foo":"bar","uid":"OLDUID"}`),
		)
		err := client.CreateBucket(ctx, "my-new-bucket", []byte(`{"foo":"bar","uid":"NEWUID"}`),
			bucketclient.CreateBucketMakeIdempotent)
		Expect(err).To(HaveOccurred())
		bcErr, ok := err.(*bucketclient.BucketClientError)
		Expect(ok).To(BeTrue())
		Expect(bcErr.StatusCode).To(Equal(http.StatusConflict))
	})

	It("returns 409 Conflict with MakeIdempotent option if bucket exists without an \"uid\" attribute", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"POST", "/default/bucket/my-new-bucket",
			httpmock.NewStringResponder(http.StatusConflict, ""),
		)
		httpmock.RegisterResponder(
			"GET", "/default/attributes/my-new-bucket",
			httpmock.NewStringResponder(200, `{"foo":"bar"}`),
		)
		err := client.CreateBucket(ctx, "my-new-bucket", []byte(`{"foo":"bar"}`),
			bucketclient.CreateBucketMakeIdempotent)
		Expect(err).To(HaveOccurred())
		bcErr, ok := err.(*bucketclient.BucketClientError)
		Expect(ok).To(BeTrue())
		Expect(bcErr.StatusCode).To(Equal(http.StatusConflict))
	})

	It("returns 409 Conflict with MakeIdempotent option if bucket exists with invalid attributes", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"POST", "/default/bucket/my-new-bucket",
			httpmock.NewStringResponder(http.StatusConflict, ""),
		)
		httpmock.RegisterResponder(
			"GET", "/default/attributes/my-new-bucket",
			httpmock.NewStringResponder(200, "NOT-JSON"),
		)
		err := client.CreateBucket(ctx, "my-new-bucket", []byte(`{"foo":"bar"}`),
			bucketclient.CreateBucketMakeIdempotent)
		Expect(err).To(HaveOccurred())
		bcErr, ok := err.(*bucketclient.BucketClientError)
		Expect(ok).To(BeTrue())
		Expect(bcErr.StatusCode).To(Equal(http.StatusConflict))
	})
})
