package bucketclient

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"fmt"
)

var _ = Describe("ListObjectVersions utility functions", func() {
	Describe("CompareVersionsListingMarkers()", func() {
		var tests = []struct {
			keyMarker1       string
			versionIdMarker1 string
			keyMarker2       string
			versionIdMarker2 string
			expectedResult   int
		}{
			{
				keyMarker1:       "",
				versionIdMarker1: "",
				keyMarker2:       "",
				versionIdMarker2: "",
				expectedResult:   0,
			},
			{
				keyMarker1:       "",
				versionIdMarker1: "",
				keyMarker2:       "foo",
				versionIdMarker2: "123",
				expectedResult:   -1,
			},
			{
				keyMarker1:       "foo",
				versionIdMarker1: "123",
				keyMarker2:       "bar",
				versionIdMarker2: "321",
				expectedResult:   1,
			},
			{
				keyMarker1:       "bar",
				versionIdMarker1: "321",
				keyMarker2:       "foo",
				versionIdMarker2: "123",
				expectedResult:   -1,
			},
			{
				keyMarker1:       "foo",
				versionIdMarker1: "321",
				keyMarker2:       "foo",
				versionIdMarker2: "123",
				expectedResult:   1,
			},
			{
				keyMarker1:       "foo",
				versionIdMarker1: "123",
				keyMarker2:       "foo",
				versionIdMarker2: "321",
				expectedResult:   -1,
			},
			{
				keyMarker1:       "foo",
				versionIdMarker1: "123",
				keyMarker2:       "foo",
				versionIdMarker2: "123",
				expectedResult:   0,
			},
		}

		for _, testCase := range tests {
			It(fmt.Sprintf("'%s:%s' <=> '%s:%s' = %d",
				testCase.keyMarker1, testCase.versionIdMarker1,
				testCase.keyMarker2, testCase.versionIdMarker2,
				testCase.expectedResult,
			), func() {
				cmp := CompareVersionsListingMarkers(
					testCase.keyMarker1, testCase.versionIdMarker1,
					testCase.keyMarker2, testCase.versionIdMarker2,
				)
				Expect(cmp).To(Equal(testCase.expectedResult))
			})
		}
	})

	Describe("truncateListObjectVersionsResponse()", func() {
		var tests = []struct {
			description         string
			listResponse        ListObjectVersionsResponse
			lastKeyMarker       string
			lastVersionIdMarker string
			expectedTruncation  *ListObjectVersionsResponse
		}{
			{
				description: "empty listing",
				listResponse: ListObjectVersionsResponse{
					Versions:    []ListObjectVersionsEntry{},
					IsTruncated: false,
				},
				lastKeyMarker:       "foo",
				lastVersionIdMarker: "123",
				expectedTruncation:  nil,
			},
			{
				description: "one entry, IsTruncated=false, does not get truncated",
				listResponse: ListObjectVersionsResponse{
					Versions: []ListObjectVersionsEntry{
						ListObjectVersionsEntry{
							Key:       "abc",
							VersionId: "123",
						},
					},
					IsTruncated: false,
				},
				lastKeyMarker:       "foo",
				lastVersionIdMarker: "123",
				expectedTruncation:  nil,
			},
			{
				description: "one entry, IsTruncated=false, with key beyond lastMarker gets truncated",
				listResponse: ListObjectVersionsResponse{
					Versions: []ListObjectVersionsEntry{
						ListObjectVersionsEntry{
							Key:       "abc",
							VersionId: "123",
						},
					},
					IsTruncated: false,
				},
				lastKeyMarker:       "abc",
				lastVersionIdMarker: "122",
				expectedTruncation: &ListObjectVersionsResponse{
					Versions:    []ListObjectVersionsEntry{},
					IsTruncated: false,
				},
			},
			{
				description: "one entry, IsTruncated=false, with key equal to lastMarker does not get truncated",
				listResponse: ListObjectVersionsResponse{
					Versions: []ListObjectVersionsEntry{
						ListObjectVersionsEntry{
							Key:       "abc",
							VersionId: "123",
						},
					},
					IsTruncated: false,
				},
				lastKeyMarker:       "abc",
				lastVersionIdMarker: "123",
				expectedTruncation:  nil,
			},
			{
				description: "one entry, IsTruncated=true, with NextMarker strictly lower than lastMarker does not get truncated",
				listResponse: ListObjectVersionsResponse{
					Versions: []ListObjectVersionsEntry{
						ListObjectVersionsEntry{
							Key:       "abc",
							VersionId: "123",
						},
					},
					IsTruncated:         true,
					NextKeyMarker:       "abc",
					NextVersionIdMarker: "123",
				},
				lastKeyMarker:       "def",
				lastVersionIdMarker: "123",
				expectedTruncation:  nil,
			},
			{
				description: "one entry, IsTruncated=true, with NextMarker equal to lastMarker does not get truncated but IsTruncated is set to false",
				listResponse: ListObjectVersionsResponse{
					Versions: []ListObjectVersionsEntry{
						ListObjectVersionsEntry{
							Key:       "abc",
							VersionId: "123",
						},
					},
					IsTruncated:         true,
					NextKeyMarker:       "abc",
					NextVersionIdMarker: "123",
				},
				lastKeyMarker:       "abc",
				lastVersionIdMarker: "123",
				expectedTruncation: &ListObjectVersionsResponse{
					Versions: []ListObjectVersionsEntry{
						ListObjectVersionsEntry{
							Key:       "abc",
							VersionId: "123",
						},
					},
					IsTruncated: false,
				},
			},
			{
				description: "one entry, IsTruncated=true, with NextMarker strictly higher than lastMarker gets truncated and IsTruncated is set to false",
				listResponse: ListObjectVersionsResponse{
					Versions: []ListObjectVersionsEntry{
						ListObjectVersionsEntry{
							Key:       "abc",
							VersionId: "123",
						},
					},
					IsTruncated:         true,
					NextKeyMarker:       "abc",
					NextVersionIdMarker: "123",
				},
				lastKeyMarker:       "aaa",
				lastVersionIdMarker: "123",
				expectedTruncation: &ListObjectVersionsResponse{
					Versions:    []ListObjectVersionsEntry{},
					IsTruncated: false,
				},
			},
			{
				description: "three entries, with no truncation and IsTruncated=true",
				listResponse: ListObjectVersionsResponse{
					Versions: []ListObjectVersionsEntry{
						ListObjectVersionsEntry{
							Key:       "abc",
							VersionId: "123",
						},
						ListObjectVersionsEntry{
							Key:       "bcd",
							VersionId: "234",
						},
						ListObjectVersionsEntry{
							Key:       "cde",
							VersionId: "345",
						},
					},
					IsTruncated:         true,
					NextKeyMarker:       "cde",
					NextVersionIdMarker: "345",
				},
				lastKeyMarker:       "cde",
				lastVersionIdMarker: "346",
				expectedTruncation:  nil,
			},
			{
				description: "three entries, of which one gets truncated, with IsTruncated set to false",
				listResponse: ListObjectVersionsResponse{
					Versions: []ListObjectVersionsEntry{
						ListObjectVersionsEntry{
							Key:       "abc",
							VersionId: "123",
						},
						ListObjectVersionsEntry{
							Key:       "bcd",
							VersionId: "234",
						},
						ListObjectVersionsEntry{
							Key:       "cde",
							VersionId: "345",
						},
					},
					IsTruncated:         true,
					NextKeyMarker:       "cde",
					NextVersionIdMarker: "345",
				},
				lastKeyMarker:       "caa",
				lastVersionIdMarker: "456",
				expectedTruncation: &ListObjectVersionsResponse{
					Versions: []ListObjectVersionsEntry{
						ListObjectVersionsEntry{
							Key:       "abc",
							VersionId: "123",
						},
						ListObjectVersionsEntry{
							Key:       "bcd",
							VersionId: "234",
						},
					},
					IsTruncated: false,
				},
			},
		}
		for _, testCase := range tests {
			It(testCase.description, func() {
				listResponse := testCase.listResponse

				truncateListObjectVersionsResponse(&listResponse,
					testCase.lastKeyMarker, testCase.lastVersionIdMarker)

				if testCase.expectedTruncation != nil {
					Expect(listResponse).To(Equal(*testCase.expectedTruncation))
				} else {
					Expect(listResponse).To(Equal(testCase.listResponse))
				}
			})
		}
	})
})
