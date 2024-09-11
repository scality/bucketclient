package bucketclient_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/scality/bucketclient/go"
)

var _ = Describe("BucketClientError", func() {
	It("Error(): with non-200 HTTP status", func() {
		myError := &bucketclient.BucketClientError{
			ApiMethod:  "SomeMethod",
			HttpMethod: "GET",
			Endpoint:   "http://localhost:9000",
			Resource:   "/some/resource",
			StatusCode: 404,
			ErrorType:  "ResourceNotFound",
			Err:        nil,
		}
		errStr := myError.Error()
		Expect(errStr).To(Equal("error in SomeMethod [GET http://localhost:9000/some/resource]: " +
			"bucketd returned HTTP status 404 ResourceNotFound"))
	})
	It("Error(): generic error", func() {
		myError := &bucketclient.BucketClientError{
			ApiMethod:  "SomeMethod",
			HttpMethod: "GET",
			Endpoint:   "http://localhost:9000",
			Resource:   "/some/resource",
			StatusCode: 0,
			ErrorType:  "",
			Err:        errors.New("OOPS"),
		}
		errStr := myError.Error()
		Expect(errStr).To(Equal("error in SomeMethod [GET http://localhost:9000/some/resource]: " +
			"HTTP request to bucketd failed: OOPS"))
		Expect(myError.Unwrap().Error()).To(Equal("OOPS"))
	})
	It("ErrorMalformedResponse() creates a specific error for malformed responses", func() {
		myError := bucketclient.ErrorMalformedResponse(
			"SomeMethod", "GET", "http://localhost:9000", "/some/resource",
			errors.New("OOPS"),
		)
		errStr := myError.Error()
		Expect(errStr).To(Equal("error in SomeMethod [GET http://localhost:9000/some/resource]: " +
			"HTTP request to bucketd failed: bucketd returned a malformed response body: " +
			"OOPS"))
	})
})
