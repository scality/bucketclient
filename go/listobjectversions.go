package bucketclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

type ListObjectVersionsOption func(*listObjectVersionsOptionSet) error

func ListObjectVersionsMarkerOption(keyMarker string, versionIdMarker string) ListObjectVersionsOption {
	return func(opts *listObjectVersionsOptionSet) error {
		opts.keyMarker = &keyMarker
		opts.versionIdMarker = &versionIdMarker
		return nil
	}
}

func ListObjectVersionsMaxKeysOption(maxKeys int) ListObjectVersionsOption {
	return func(opts *listObjectVersionsOptionSet) error {
		if maxKeys < 0 || maxKeys > 10000 {
			return fmt.Errorf("maxKeys=%d is out of the valid range [0, 10000]", maxKeys)
		}
		opts.maxKeys = &maxKeys
		return nil
	}
}

// ListObjectVersionsLastMarkerOption option makes the listing behave
// as if the bucket contains no object which key/versionId is strictly
// higher than the pair "lastKeyMarker/lastVersionIdMarker".
//
// Note: this option is not implemented natively by bucketd, hence the
// Go client may truncate the result and adjust the "IsTruncated"
// field accordingly, before returning the truncated response to the
// client.
func ListObjectVersionsLastMarkerOption(lastKeyMarker string, lastVersionIdMarker string) ListObjectVersionsOption {
	return func(opts *listObjectVersionsOptionSet) error {
		opts.lastKeyMarker = &lastKeyMarker
		opts.lastVersionIdMarker = &lastVersionIdMarker
		return nil
	}
}

type ListObjectVersionsEntry struct {
	Key       string `json:"key"`
	VersionId string `json:"versionId"`
	Value     string `json:"value"`
}

type ListObjectVersionsResponse struct {
	Versions            []ListObjectVersionsEntry
	CommonPrefixes      []string
	IsTruncated         bool
	NextKeyMarker       string `json:",omitempty"`
	NextVersionIdMarker string `json:",omitempty"`
}

type listObjectVersionsOptionSet struct {
	keyMarker           *string
	versionIdMarker     *string
	maxKeys             *int
	lastKeyMarker       *string
	lastVersionIdMarker *string
}

func parseListObjectVersionsOptions(opts []ListObjectVersionsOption) (listObjectVersionsOptionSet, error) {
	parsedOpts := listObjectVersionsOptionSet{}
	for _, opt := range opts {
		err := opt(&parsedOpts)
		if err != nil {
			return parsedOpts, err
		}
	}
	return parsedOpts, nil
}

func (client *BucketClient) ListObjectVersions(ctx context.Context,
	bucketName string, opts ...ListObjectVersionsOption) (*ListObjectVersionsResponse, error) {
	resource := fmt.Sprintf("/default/bucket/%s", bucketName)
	query := url.Values{}
	query.Set("listingType", "DelimiterVersions")

	options, err := parseListObjectVersionsOptions(opts)
	if err != nil {
		return nil, &BucketClientError{
			"ListObjectVersions", "GET", client.Endpoint, resource, 0, "", err,
		}
	}
	if options.keyMarker != nil {
		query.Set("keyMarker", *options.keyMarker)
		query.Set("versionIdMarker", *options.versionIdMarker)
	}
	if options.maxKeys != nil {
		query.Set("maxKeys", strconv.Itoa(*options.maxKeys))
	}
	u, _ := url.Parse(resource)
	u.RawQuery = query.Encode()
	resource = u.String()
	responseBody, err := client.Request(ctx, "ListObjectVersions", "GET", resource)
	if err != nil {
		return nil, err
	}
	var parsedResponse ListObjectVersionsResponse
	jsonErr := json.Unmarshal(responseBody, &parsedResponse)
	if jsonErr != nil {
		return nil, ErrorMalformedResponse("ListObjectVersions", "GET",
			client.Endpoint, resource, jsonErr)
	}
	if options.lastKeyMarker != nil {
		truncateListObjectVersionsResponse(&parsedResponse,
			*options.lastKeyMarker, *options.lastVersionIdMarker)
	}
	return &parsedResponse, nil
}
