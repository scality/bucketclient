package bucketclient_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jarcoal/httpmock"

	"github.com/scality/bucketclient/go"
)

var _ = Describe("ListBasic()", func() {
	It("returns an empty listing result", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"GET", "/default/bucket/my-bucket?listingType=Basic",
			httpmock.NewStringResponder(200, "[]"))

		Expect(client.ListBasic(ctx, "my-bucket")).To(Equal(
			&bucketclient.ListBasicResponse{}))
	})

	It("returns a non-empty listing result with no param", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"GET", "/default/bucket/my-bucket?listingType=Basic",
			httpmock.NewStringResponder(200, `[
        {"key": "fop", "value": "fopvalue"},
        {"key": "goo", "value": "goovalue"},
        {"key": "hop", "value": "hopvalue"}
]
`))

		Expect(client.ListBasic(ctx, "my-bucket")).To(Equal(&bucketclient.ListBasicResponse{
			bucketclient.ListBasicEntry{Key: "fop", Value: "fopvalue"},
			bucketclient.ListBasicEntry{Key: "goo", Value: "goovalue"},
			bucketclient.ListBasicEntry{Key: "hop", Value: "hopvalue"},
		}))
	})

	It("returns a non-empty listing result with URL-encoded gt, lt, maxKeys params, without values", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"GET", "/default/bucket/my-bucket?gt=eee%2F123&listingType=Basic&lt=ijk+456&maxKeys=3&values=false",
			httpmock.NewStringResponder(200, `[
        {"key": "fop"},
        {"key": "goo"},
        {"key": "hop"}
]
`))

		Expect(client.ListBasic(ctx, "my-bucket",
			bucketclient.ListBasicGTOption("eee/123"),
			bucketclient.ListBasicLTOption("ijk 456"),
			bucketclient.ListBasicMaxKeysOption(3),
			bucketclient.ListBasicNoValuesOption(),
		)).To(Equal(&bucketclient.ListBasicResponse{
			bucketclient.ListBasicEntry{Key: "fop"},
			bucketclient.ListBasicEntry{Key: "goo"},
			bucketclient.ListBasicEntry{Key: "hop"},
		}))
	})

	It("returns an error with malformed response", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"GET", "/default/bucket/my-bucket?listingType=Basic",
			httpmock.NewStringResponder(200, "[OOPS"),
		)

		_, err := client.ListBasic(ctx, "my-bucket")
		Expect(err).To(MatchError(ContainSubstring("malformed response body")))
	})

})
