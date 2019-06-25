package sensorsapi
type ConsumerType int
const (
	_ ConsumerType =iota
	Logging 
	ConcurrentLogging
	Debug
	Default
	Batch
)
