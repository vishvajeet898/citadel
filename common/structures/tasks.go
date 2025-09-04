package structures

import "time"

type TaskIdPreviousStateStruct struct {
	TaskId         uint
	PreviousStatus string
}

type VisitDetailsForTask struct {
	VisitId           string     `json:"visit_id"`
	SampleCollectedAt *time.Time `json:"sample_collected_at"`
	SampleReceivedAt  *time.Time `json:"sample_received_at"`
}
