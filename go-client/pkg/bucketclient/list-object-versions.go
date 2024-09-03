package bucketclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type ListObjectVersionsEntry struct {
	Key string `json:"key"`
	VersionId string `json:"versionId"`
	Value string `json:"value"`
}

type ListObjectVersionsResponse struct {
	Versions []ListObjectVersionsEntry
	CommonPrefixes []string
	IsTruncated bool
	NextKeyMarker string
	NextVersionIdMarker string
}

func ListObjectVersions(bucketdUrl string, bucketName string,
	keyMarker string, versionIdMarker string, maxKeys int) (*ListObjectVersionsResponse, error) {
	listObjectVersionsURL := fmt.Sprintf(
		"%s/default/bucket/%s?listingType=DelimiterVersions&keyMarker=%s&versionIdMarker=%s&maxKeys=%d",
		bucketdUrl, bucketName,
		url.QueryEscape(keyMarker), url.QueryEscape(versionIdMarker), maxKeys)
	listResp, err := http.Get(listObjectVersionsURL)
	if err != nil {
		return nil, err
	}
	defer listResp.Body.Close()

	if listResp.StatusCode != 200 {
		return nil, fmt.Errorf("bucketd returned HTTP status %s",
			listResp.Status)
	}
	listBody, err := io.ReadAll(listResp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	var parsedResponse = new(ListObjectVersionsResponse)
	err = json.Unmarshal(listBody, parsedResponse)
	if err != nil {
		return nil, fmt.Errorf("bucketd returned a malformed response body: %w", err)
	}

	return parsedResponse, nil
}
