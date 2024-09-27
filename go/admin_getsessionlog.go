package bucketclient

import (
	"context"
	"encoding/json"
	"fmt"
)

// AdminGetSessionLog returns a range of Raft oplog from the given raft session.
// sessionId is the raft session ID
// beginSeq is the first raft sequence number to fetch
// nRecords is the maximum number of records to fetch starting from beginSeq
// if targetLeader is true, it fetches the oplog from the leader,
// otherwise fetches from one of the followers
func (client *BucketClient) AdminGetSessionLog(ctx context.Context,
	sessionId int, beginSeq int64, nRecords int, targetLeader bool) (*SessionLogResponse, error) {
	resource := fmt.Sprintf("/_/raft_sessions/%d/log?begin=%d&limit=%d", sessionId, beginSeq, nRecords)
	if targetLeader {
		resource += "&target_leader=true"
	}
	responseBody, err := client.Request(ctx, "AdminGetSessionLog", "GET", resource)
	if err != nil {
		return nil, err
	}
	var parsedResponse SessionLogResponse
	jsonErr := json.Unmarshal(responseBody, &parsedResponse)
	if jsonErr != nil {
		return nil, ErrorMalformedResponse("AdminGetSessionsLog",
			"GET", client.Endpoint, resource, jsonErr)
	}
	return &parsedResponse, nil
}
