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

const (
	MetricsNamespace = "s3_metadata_bucketclient"
)

var MetricsSummaryDefaultObjectives = map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001, 1.0: 0.0001}
