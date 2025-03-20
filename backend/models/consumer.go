package models

type ConsumerGroupInfo struct {
	GroupID     string         `json:"group_id"`
	Members     int            `json:"members"`
	Topics      int            `json:"topics"`
	Lag         int64          `json:"lag"`
	Coordinator int32          `json:"coordinator"`
	State       string         `json:"state"`
	TopicLags   []TopicLagInfo `json:"topic_lags,omitempty"`
}

type TopicLagInfo struct {
	Topic     string `json:"topic"`
	Partition int32  `json:"partition"`
	Lag       int64  `json:"lag"`
}

type ConsumerCreateRequest struct {
	Topic string `json:"topic" binding:"required"`
	Group string `json:"group" binding:"required"`
	ID    string `json:"id" binding:"required"`
}

type ConsumerCreateResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Poll   string `json:"poll"`
}

type MessageProduceRequest struct {
	Value interface{} `json:"value" binding:"required"`
	Key   string      `json:"key,omitempty"`
}


