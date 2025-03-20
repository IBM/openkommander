package models

type BrokerInfo struct {
	ID               int32  `json:"id"`
	Host             string `json:"host"`
	Port             int32  `json:"port"`
	PartitionsLeader int    `json:"partitions_leader"`
	Partitions       int    `json:"partitions"`
	InSyncPartitions int    `json:"in_sync_partitions"`
}


