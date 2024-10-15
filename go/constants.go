package bucketclient

type DBMethodType int

const (
	DBMethodCreate        DBMethodType = 0
	DBMethodDelete        DBMethodType = 1
	DBMethodGet           DBMethodType = 2
	DBMethodPut           DBMethodType = 3
	DBMethodList          DBMethodType = 4
	DBMethodDel           DBMethodType = 5
	DBMethodGetAttributes DBMethodType = 6
	DBMethodPutAttributes DBMethodType = 7
	DBMethodBatch         DBMethodType = 8
	DBMethodNoop          DBMethodType = 9
)
