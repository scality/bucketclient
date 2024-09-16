package bucketclient

type MemberInfo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	AdminPort   int    `json:"adminPort"`
	MDClusterId string `json:"mdClusterId"`
}

type SessionInfo struct {
	ID                int          `json:"id"`
	RaftMembers       []MemberInfo `json:"raftMembers"`
	ConnectedToLeader bool         `json:"connectedToLeader"`
}
