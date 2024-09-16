package bucketclient

import (
	"context"
	"encoding/json"
	"fmt"
)

// AdminGetAllSessionsInfo returns raft session info for all Metadata
// raft sessions available on the S3C deployment.
func (client *BucketClient) AdminGetAllSessionsInfo(ctx context.Context) ([]SessionInfo, error) {
	responseBody, err := client.Request(ctx, "AdminGetAllSessionsInfo", "GET", "/_/raft_sessions")
	if err != nil {
		return nil, err
	}

	var parsedInfo []SessionInfo
	jsonErr := json.Unmarshal(responseBody, &parsedInfo)
	if jsonErr != nil {
		return nil, ErrorMalformedResponse("AdminGetAllSessionsInfo",
			"GET", client.Endpoint, "/_/raft_sessions", jsonErr)
	}
	return parsedInfo, nil
}

// AdminGetSessionInfo returns raft session info for the given raft session ID.
// Returns nil and an error if the raft session doesn't exist, or if a request
// error occurs.
func (client *BucketClient) AdminGetSessionInfo(ctx context.Context,
	sessionId int) (*SessionInfo, error) {
	// When querying /_/raft_sessions/X/info, bucketd returns a
	// status 500 if the raft session doesn't exist, which is hard
	// to distinguish from a generic type of failure. For this
	// reason, instead, we fetch the info for all raft sessions
	// then lookup the one we want.
	sessionsInfo, err := client.AdminGetAllSessionsInfo(ctx)
	if err != nil {
		return nil, err
	}
	for _, sessionInfo := range sessionsInfo {
		if sessionInfo.ID == sessionId {
			return &sessionInfo, nil
		}
	}
	// raft session does not exist: return a 404 status as if coming from bucketd
	return nil, &BucketClientError{
		"AdminGetSessionInfo", "GET", client.Endpoint, "/_/raft_sessions",
		404, "RaftSessionNotFound",
		fmt.Errorf("no such raft session: %d", sessionId),
	}
}
