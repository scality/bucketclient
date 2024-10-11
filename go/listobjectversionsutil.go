package bucketclient

import ()

// CompareVersionsListingMarkers is a helper function that returns -1,
// 0, or 1 if the pair keyMarker1/versionIdMarker1 is
// lexicographically, respectively strictly lower, equal, or strictly
// higher than the pair keyMarker2/versionIdMarker2.
func CompareVersionsListingMarkers(keyMarker1 string, versionIdMarker1 string,
	keyMarker2 string, versionIdMarker2 string) int {
	// if key markers are different, versionId markers are ignored
	if keyMarker1 != keyMarker2 {
		if keyMarker1 < keyMarker2 {
			return -1
		}
		return 1
	}
	// if key markers are equal, versionId markers are compared
	if versionIdMarker1 != versionIdMarker2 {
		if versionIdMarker1 < versionIdMarker2 {
			return -1
		}
		return 1
	}
	return 0
}

// truncateListObjectVersionsResponse discards entries which
// key/versionId pair is strictly higher than the
// lastKeyMarker/lastVersionIdMarker pair, and may also change the
// IsTruncated attribute.
func truncateListObjectVersionsResponse(listResponse *ListObjectVersionsResponse,
	lastKeyMarker string, lastVersionIdMarker string) {
	if listResponse.IsTruncated {
		cmp := CompareVersionsListingMarkers(
			listResponse.NextKeyMarker, listResponse.NextVersionIdMarker,
			lastKeyMarker, lastVersionIdMarker)
		if cmp < 0 {
			return
		}
	}
	var i int
	for i = len(listResponse.Versions) - 1; i >= 0; i -= 1 {
		cmp := CompareVersionsListingMarkers(
			listResponse.Versions[i].Key, listResponse.Versions[i].VersionId,
			lastKeyMarker, lastVersionIdMarker)
		if cmp <= 0 {
			break
		}
	}
	listResponse.IsTruncated = false
	listResponse.NextKeyMarker = ""
	listResponse.NextVersionIdMarker = ""
	if i+1 < len(listResponse.Versions) {
		listResponse.Versions = listResponse.Versions[0 : i+1]
	}
}
