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

type SessionLogResponse struct {
	Info SessionLogInfo     `json:"info"`
	Log  []SessionLogRecord `json:"log"`
}

type SessionLogInfo struct {
	Start int64 `json:"start"`
	CSeq  int64 `json:"cseq"`
	Prune int64 `json:"prune"`
}

type SessionLogRecord struct {
	Bucket    string            `json:"db"`
	DBMethod  DBMethodType      `json:"method"`
	Timestamp string            `json:"timestamp,omitempty"`
	Entries   []SessionLogEntry `json:"entries"`
}

type SessionLogEntry struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
	Type  string `json:"type,omitempty"`
}
