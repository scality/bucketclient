package bucketclient

type MetastoreEntry struct {
	Name          string `json:"name"`
	Attributes    string `json:"attributes"`
	Creating      bool   `json:"creating"`
	Deleting      bool   `json:"deleting"`
	ID            string `json:"id"`
	RaftSessionID int    `json:"raftSessionID"`
	Version       int    `json:"version"`
	RaftSession   string `json:"raftSession"`
}
