package types

// SentinelMasterCmdResult represents the output of the "sentinel master" command
type SentinelMasterCmdResult struct {
	Name                  string `redis:"name"`
	IP                    string `redis:"ip"`
	Port                  int    `redis:"port"`
	RunID                 string `redis:"runid"`
	Flags                 string `redis:"flags"`
	LinkPendingCommands   int    `redis:"link-pending-commands"`
	LinkRefcount          int    `redis:"link-refcount"`
	LastPingSet           int    `redis:"last-ping-sent"`
	LastOkPingReply       int    `redis:"last-ok-ping-reply"`
	LastPingReply         int    `redis:"last-ping-reply"`
	DownAfterMilliseconds int    `redis:"down-after-milliseconds"`
	InfoRefresh           int    `redis:"info-refresh"`
	RoleReported          string `redis:"role-reported"`
	RoleReportedTime      int    `redis:"role-reported-time"`
	ConfigEpoch           int    `redis:"config-epoch"`
	NumSlaves             int    `redis:"num-slaves"`
	NumOtherSentinels     int    `redis:"num-other-sentinels"`
	Quorum                int    `redis:"quorum"`
	FailoverTimeout       int    `redis:"failover-timeout"`
	ParallelSyncs         int    `redis:"parallel-syncs"`
}

// SentinelSlaveCmdResult represents the output of the "sentinel slave" command
type SentinelSlaveCmdResult struct {
	Name                  string `redis:"name"`
	IP                    string `redis:"ip"`
	Port                  int    `redis:"port"`
	RunID                 string `redis:"runid"`
	Flags                 string `redis:"flags"`
	LinkPendingCommands   int    `redis:"link-pending-commands"`
	LinkRefcount          int    `redis:"link-refcount"`
	LastPingSet           int    `redis:"last-ping-sent"`
	LastOkPingReply       int    `redis:"last-ok-ping-reply"`
	LastPingReply         int    `redis:"last-ping-reply"`
	DownAfterMilliseconds int    `redis:"down-after-milliseconds"`
	InfoRefresh           int    `redis:"info-refresh"`
	RoleReported          string `redis:"role-reported"`
	RoleReportedTime      int    `redis:"role-reported-time"`
	MasterLinkDownTime    int    `redis:"master-link-down-time"`
	MasterLinkStatus      string `redis:"master-link-status"`
	MasterHost            string `redis:"master-host"`
	MasterPort            string `redis:"master-port"`
	SlavePriority         string `redis:"slave-priority"`
	SlaveReplOffset       int    `redis:"slave-repl-offset"`
}
