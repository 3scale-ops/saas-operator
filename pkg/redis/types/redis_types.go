package types

// Role represents the role of a redis server within a shard
type Role string

const (
	// Master is the master role in a shard. Under normal circumstances, only
	// a server in the shard can be master at a given time
	Master Role = "master"
	// Slave are servers within the shard that replicate data from the master
	// for data high availabilty purposes
	Slave Role = "slave"
	// Unknown represents a state in which the role of the server is still unknown
	Unknown Role = "unknown"
)
