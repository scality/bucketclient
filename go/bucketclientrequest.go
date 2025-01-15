package bucketclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type requestOptionSet struct {
	requestBody            []byte
	requestBodyContentType string
	requestUIDs            string
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

func RequestUIDsOption(uids string) RequestOption {
	return func(ros *requestOptionSet) {
		ros.requestUIDs = uids
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

// updateMetrics is a local helper to update Prometheus metrics in a generic way before
// returning a response or an error to the API caller.
func (client *BucketClient) updateMetrics(apiMethod, httpMethod string,
	startTime time.Time, requestBody []byte,
	httpCode int, responseBody []byte) {
	if client.Metrics == nil {
		return
	}
	labels := prometheus.Labels{
		"endpoint": client.Endpoint,
		"method":   httpMethod,
		"action":   apiMethod,
		"code":     strconv.Itoa(httpCode),
	}
	elapsedSeconds := time.Since(startTime).Seconds()

	client.Metrics.RequestsTotal.With(labels).Inc()
	client.Metrics.RequestDurationSeconds.With(labels).Observe(elapsedSeconds)

	// only update this metric when the request body exists, to avoid unneeded metrics
	if requestBody != nil {
		client.Metrics.RequestBytesSentTotal.With(labels).Add(float64(len(requestBody)))
	}

	var responseBodyLength int
	if responseBody != nil {
		responseBodyLength = len(responseBody)
	}
	client.Metrics.ResponseBytesReceivedTotal.With(labels).Add(float64(responseBodyLength))
}

func (client *BucketClient) Request(ctx context.Context,
	apiMethod string, httpMethod string, resource string, opts ...RequestOption) ([]byte, error) {
	var response *http.Response
	var err error

	startTime := time.Now()
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

			if options.requestUIDs != "" {
				request.Header.Add("x-scal-request-uids", options.requestUIDs)
			}

			if options.idempotent {
				request.Header["Idempotency-Key"] = []string{}
			}
			response, err = http.DefaultClient.Do(request)
		}
	}
	if err != nil {
		client.updateMetrics(apiMethod, httpMethod, startTime, options.requestBody, 0, nil)
		return nil, &BucketClientError{
			apiMethod, httpMethod, client.Endpoint, resource, 0, "", err,
		}
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		// We have a HTTP status code but we couldn't read the whole response body,
		// so use "0" as status code for a generic transport error
		client.updateMetrics(apiMethod, httpMethod, startTime, options.requestBody, 0, nil)
		return nil, &BucketClientError{
			apiMethod, httpMethod, client.Endpoint, resource, 0, "",
			fmt.Errorf("error reading response body: %w", err),
		}
	}
	client.updateMetrics(apiMethod, httpMethod, startTime, options.requestBody,
		response.StatusCode, responseBody)

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
	return responseBody, nil
}
