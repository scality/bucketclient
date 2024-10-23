package bucketclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

type ListBasicOption func(*listBasicOptionSet) error

// ListBasicGTOption only lists keys greater than the given argument
func ListBasicGTOption(gt string) ListBasicOption {
	return func(opts *listBasicOptionSet) error {
		opts.gt = &gt
		return nil
	}
}

// ListBasicGTEOption only lists keys greater or equal to the given argument
func ListBasicGTEOption(gte string) ListBasicOption {
	return func(opts *listBasicOptionSet) error {
		opts.gte = &gte
		return nil
	}
}

// ListBasicLTOption only lists keys less than the given argument
func ListBasicLTOption(lt string) ListBasicOption {
	return func(opts *listBasicOptionSet) error {
		opts.lt = &lt
		return nil
	}
}

// ListBasicLTEOption only lists keys less or equal to the given argument
func ListBasicLTEOption(lte string) ListBasicOption {
	return func(opts *listBasicOptionSet) error {
		opts.lte = &lte
		return nil
	}
}

// ListBasicMaxKeysOption limits the number of returned keys (default and maximum is 10000).
func ListBasicMaxKeysOption(maxKeys int) ListBasicOption {
	return func(opts *listBasicOptionSet) error {
		if maxKeys < 0 || maxKeys > 10000 {
			return fmt.Errorf("maxKeys=%d is out of the valid range [0, 10000]", maxKeys)
		}
		opts.maxKeys = &maxKeys
		return nil
	}
}

// ListBasicNoKeysOption declares that keys are not needed in the
// result entries and may be returned empty.
//
// Note: keys may still be returned until ARSN-438 is fixed.
func ListBasicNoKeysOption() ListBasicOption {
	return func(opts *listBasicOptionSet) error {
		opts.noKeys = true
		return nil
	}
}

// ListBasicNoValuesOption declares that values are not needed in the
// result entries and may be returned empty.
//
// Note: values may still be returned until ARSN-438 is fixed.
func ListBasicNoValuesOption() ListBasicOption {
	return func(opts *listBasicOptionSet) error {
		opts.noValues = true
		return nil
	}
}

type ListBasicEntry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ListBasicResponse []ListBasicEntry

type listBasicOptionSet struct {
	gt       *string
	gte      *string
	lt       *string
	lte      *string
	maxKeys  *int
	noKeys   bool
	noValues bool
}

func parseListBasicOptions(opts []ListBasicOption) (listBasicOptionSet, error) {
	parsedOpts := listBasicOptionSet{}
	for _, opt := range opts {
		err := opt(&parsedOpts)
		if err != nil {
			return parsedOpts, err
		}
	}
	return parsedOpts, nil
}

func (client *BucketClient) ListBasic(ctx context.Context,
	bucketName string, opts ...ListBasicOption) (*ListBasicResponse, error) {
	resource := fmt.Sprintf("/default/bucket/%s", bucketName)
	query := url.Values{}
	query.Set("listingType", "Basic")

	options, err := parseListBasicOptions(opts)
	if err != nil {
		return nil, &BucketClientError{
			"ListBasic", "GET", client.Endpoint, resource, 0, "", err,
		}
	}
	if options.gt != nil {
		query.Set("gt", *options.gt)
	}
	if options.gte != nil {
		query.Set("gte", *options.gte)
	}
	if options.lt != nil {
		query.Set("lt", *options.lt)
	}
	if options.lte != nil {
		query.Set("lte", *options.lte)
	}
	if options.maxKeys != nil {
		query.Set("maxKeys", strconv.Itoa(*options.maxKeys))
	}
	if options.noKeys {
		query.Set("keys", "false")
	}
	if options.noValues {
		query.Set("values", "false")
	}
	u, _ := url.Parse(resource)
	u.RawQuery = query.Encode()
	resource = u.String()
	responseBody, err := client.Request(ctx, "ListBasic", "GET", resource)
	if err != nil {
		return nil, err
	}
	var parsedResponse ListBasicResponse
	jsonErr := json.Unmarshal(responseBody, &parsedResponse)
	if jsonErr != nil {
		return nil, ErrorMalformedResponse("ListBasic", "GET",
			client.Endpoint, resource, jsonErr)
	}
	return &parsedResponse, nil
}
