package bucketclient_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"io"
	"net/http"
	"strings"

	"github.com/jarcoal/httpmock"
	"github.com/prometheus/client_golang/prometheus"
	promTestutil "github.com/prometheus/client_golang/prometheus/testutil"

	"github.com/scality/bucketclient/go"
)

var _ = Describe("BucketClientMetrics", func() {
	Describe("RegisterMetrics()", func() {
		It("enables metrics collection for all requests of all clients with metrics enabled", func(ctx SpecContext) {
			client1 := bucketclient.New("http://localhost:9000")
			Expect(client1).ToNot(BeNil())
			client2 := bucketclient.New("http://localhost:9001")
			Expect(client2).ToNot(BeNil())

			registry := prometheus.NewPedanticRegistry()
			client1.EnableMetrics(registry)
			client2.EnableMetrics(registry)

			mockAttributes := `{"foo":"bar"}`
			httpmock.RegisterResponder(
				"GET", "http://localhost:9000/default/attributes/my-bucket",
				httpmock.NewStringResponder(200, mockAttributes),
			)
			expectedBatch := `{"batch":[{"key":"foo","value":"{}"}]}`
			batchErrorResponse := "OOPS!"
			httpmock.RegisterResponder(
				"POST", "http://localhost:9001/default/batch/my-bucket",
				func(req *http.Request) (*http.Response, error) {
					defer req.Body.Close()
					Expect(io.ReadAll(req.Body)).To(Equal([]byte(expectedBatch)))
					return httpmock.NewStringResponse(500, batchErrorResponse), nil
				},
			)

			_, err := client1.GetBucketAttributes(ctx, "my-bucket")
			Expect(err).ToNot(HaveOccurred())

			err = client2.PostBatch(ctx, "my-bucket", []bucketclient.PostBatchEntry{
				{Key: "foo", Value: "{}"},
			})
			Expect(err).To(HaveOccurred())

			Expect(promTestutil.GatherAndCompare(registry, strings.NewReader(`
# HELP s3_metadata_bucketclient_requests_total Number of requests processed
# TYPE s3_metadata_bucketclient_requests_total counter
s3_metadata_bucketclient_requests_total{action="GetBucketAttributes",code="200",endpoint="http://localhost:9000",method="GET"} 1
s3_metadata_bucketclient_requests_total{action="PostBatch",code="500",endpoint="http://localhost:9001",method="POST"} 1

# HELP s3_metadata_bucketclient_request_bytes_sent_total Number of request body bytes sent to bucketd
# TYPE s3_metadata_bucketclient_request_bytes_sent_total counter
s3_metadata_bucketclient_request_bytes_sent_total{action="PostBatch",code="500",endpoint="http://localhost:9001",method="POST"} 38

# HELP s3_metadata_bucketclient_response_bytes_received_total Number of response body bytes received from bucketd
# TYPE s3_metadata_bucketclient_response_bytes_received_total counter
s3_metadata_bucketclient_response_bytes_received_total{action="GetBucketAttributes",code="200",endpoint="http://localhost:9000",method="GET"} 13
s3_metadata_bucketclient_response_bytes_received_total{action="PostBatch",code="500",endpoint="http://localhost:9001",method="POST"} 5
`),
				"s3_metadata_bucketclient_requests_total",
				"s3_metadata_bucketclient_request_bytes_sent_total",
				"s3_metadata_bucketclient_response_bytes_received_total",
			)).To(Succeed())
		})
	})
})
