package bucketclient

import (
	"context"
	"fmt"
)

// AdminGetSessionDBBackend returns the database backend type of the given RAFT session.
// Currently, it returns either "leveldb" or "rocksdb".
// Returns an empty string and an error if the RAFT session doesn't exist, or if a request
// error occurs.
func (client *BucketClient) AdminGetSessionDBBackend(ctx context.Context, sessionId int,
	opts ...RequestOption) (string, error) {
	resource := fmt.Sprintf("/_/raft_sessions/%d/dbBackend", sessionId)
	responseBody, err := client.Request(ctx, "AdminGetSessionDBBackend", "GET", resource, opts...)
	if err != nil {
		bcErr := err.(*BucketClientError)
		// for backward-compatibility with older Metadata versions: if the route
		// doesn't exist, assume the database type is "leveldb"
		if bcErr.StatusCode == 400 {
			return "leveldb", nil
		}
		return "", err
	}
	dbBackend := string(responseBody)
	return dbBackend, nil
}
