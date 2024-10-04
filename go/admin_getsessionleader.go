package bucketclient

import (
	"context"
	"encoding/json"
	"fmt"
)

// AdminGetSessionLeader returns the member info for the leader of the
// given raft session ID.

// Returns nil and an error if the raft session doesn't exist, if
// bucketd is not connected to the leader, or if a request error
// occurs.
func (client *BucketClient) AdminGetSessionLeader(ctx context.Context, sessionId int) (*MemberInfo, error) {
	resource := fmt.Sprintf("/_/raft_sessions/%d/leader", sessionId)
	responseBody, err := client.Request(ctx, "AdminGetSessionLeader", "GET", resource)
	if err != nil {
		return nil, err
	}

	var parsedInfo MemberInfo
	jsonErr := json.Unmarshal(responseBody, &parsedInfo)
	if jsonErr != nil {
		return nil, ErrorMalformedResponse("AdminGetSessionLeader",
			"GET", client.Endpoint, resource, jsonErr)
	}
	return &parsedInfo, nil
}
