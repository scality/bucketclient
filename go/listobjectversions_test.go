package bucketclient_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jarcoal/httpmock"

	"github.com/scality/bucketclient/go"
)

var _ = Describe("ListObjectVersions()", func() {
	It("returns an empty listing result", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"GET", "/default/bucket/my-bucket?listingType=DelimiterVersions",
			httpmock.NewStringResponder(200, `{
    "Versions": [],
    "CommonPrefixes": [],
    "IsTruncated": false
}
`))

		Expect(client.ListObjectVersions(ctx, "my-bucket")).To(Equal(
			&bucketclient.ListObjectVersionsResponse{
				Versions:       []bucketclient.ListObjectVersionsEntry{},
				CommonPrefixes: []string{},
				IsTruncated:    false,
			}))
	})

	It("returns a non-empty listing result with URL-encoded marker, maxKeys and truncation", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"GET", "/default/bucket/my-bucket?keyMarker=foo%2Fbar&listingType=DelimiterVersions"+
				"&maxKeys=3&versionIdMarker=123+4",
			httpmock.NewStringResponder(200, `{
    "Versions": [
        {"key": "fop", "versionId": "123"},
        {"key": "goo", "versionId": "124"},
        {"key": "hop", "versionId": "125"}
    ],
    "CommonPrefixes": [],
    "IsTruncated": true,
    "NextKeyMarker": "hop",
    "NextVersionIdMarker": "125"
}
`))

		Expect(client.ListObjectVersions(ctx, "my-bucket",
			bucketclient.ListObjectVersionsMarkerOption("foo/bar", "123 4"),
			bucketclient.ListObjectVersionsMaxKeysOption(3),
			bucketclient.ListObjectVersionsLastMarkerOption("hoo", "126"),
		)).To(Equal(&bucketclient.ListObjectVersionsResponse{
			Versions: []bucketclient.ListObjectVersionsEntry{
				bucketclient.ListObjectVersionsEntry{Key: "fop", VersionId: "123"},
				bucketclient.ListObjectVersionsEntry{Key: "goo", VersionId: "124"},
			},
			CommonPrefixes: []string{},
			IsTruncated:    false,
		}))
	})

	It("returns an error with malformed response", func(ctx SpecContext) {
		httpmock.RegisterResponder(
			"GET", "/default/bucket/my-bucket?listingType=DelimiterVersions",
			httpmock.NewStringResponder(200, "{OOPS"),
		)

		_, err := client.ListObjectVersions(ctx, "my-bucket")
		Expect(err).To(MatchError(ContainSubstring("malformed response body")))
	})

})
