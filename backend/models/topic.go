package models

type TopicInfo struct {
	Name              string `json:"name"`
	Partitions        int32  `json:"partitions"`
	ReplicationFactor int16  `json:"replication_factor"`
	Internal          bool   `json:"internal"`
	Replicas          int    `json:"replicas"`
	InSyncReplicas    int    `json:"in_sync_replicas"`
	CleanupPolicy     string `json:"cleanup_policy"`
}

type TopicCreateRequest struct {
	Name              string `json:"name" binding:"required"`
	Partitions        int32  `json:"partitions" binding:"required"`
	ReplicationFactor int16  `json:"replication_factor" binding:"required"`
}

type TopicDetail struct {
	Name              string  `json:"name"`
	Partitions        int32   `json:"partitions"`
	ReplicationFactor int16   `json:"replication_factor"`
	PartitionIDs      []int32 `json:"partition_ids"`
}


