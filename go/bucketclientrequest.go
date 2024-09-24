package bucketclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type requestOptionSet struct {
	requestBody            []byte
	requestBodyContentType string
	idempotent             bool
}

type RequestOption func(*requestOptionSet)

func RequestBodyOption(body []byte) RequestOption {
	return func(ros *requestOptionSet) {
		ros.requestBody = body
	}
}

func RequestBodyContentTypeOption(contentType string) RequestOption {
	return func(ros *requestOptionSet) {
		ros.requestBodyContentType = contentType
	}
}

func RequestIdempotent(ros *requestOptionSet) {
	ros.idempotent = true
}

func parseRequestOptions(opts ...RequestOption) (requestOptionSet, error) {
	parsedOpts := requestOptionSet{
		requestBody:            nil,
		requestBodyContentType: "",
		idempotent:             false,
	}
	for _, opt := range opts {
		opt(&parsedOpts)
	}
	return parsedOpts, nil
}

func (client *BucketClient) Request(ctx context.Context,
	apiMethod string, httpMethod string, resource string, opts ...RequestOption) ([]byte, error) {
	var response *http.Response
	var err error

	options, err := parseRequestOptions(opts...)
	if err == nil {
		url := fmt.Sprintf("%s%s", client.Endpoint, resource)

		var requestBodyReader io.Reader = nil
		if options.requestBody != nil {
			requestBodyReader = bytes.NewReader(options.requestBody)
		}
		var request *http.Request
		request, err = http.NewRequestWithContext(ctx, httpMethod, url, requestBodyReader)
		if err == nil {
			if options.requestBodyContentType != "" {
				request.Header.Add("Content-Type", string(options.requestBodyContentType))
			}
			if options.idempotent {
				request.Header["Idempotency-Key"] = []string{}
			}
			response, err = http.DefaultClient.Do(request)
		}
	}
	if err != nil {
		return nil, &BucketClientError{
			apiMethod, httpMethod, client.Endpoint, resource, 0, "", err,
		}
	}
	if response.Body != nil {
		defer response.Body.Close()
	}

	if response.StatusCode/100 != 2 {
		splitStatus := strings.Split(response.Status, " ")
		errorType := ""
		if len(splitStatus) == 2 {
			errorType = splitStatus[1]
		}
		return nil, &BucketClientError{
			apiMethod, httpMethod, client.Endpoint, resource,
			response.StatusCode, errorType, nil,
		}
	}
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, &BucketClientError{
			apiMethod, httpMethod, client.Endpoint, resource, 0, "",
			fmt.Errorf("error reading response body: %w", err),
		}
	}
	return responseBody, nil
}
