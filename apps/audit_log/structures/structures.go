package structures

import (
	"time"

	commonStructures "github.com/Orange-Health/citadel/common/structures"
)

// @swagger:response AuditLog
type AuditLog struct {
	commonStructures.BaseStruct
	// The task ID of the AudiLog.
	// example: 1
	TaskID uint `json:"task_id"`
	// The log type of the AuditLog.
	// example: "Task"
	LogType string `json:"log_type"`
	// The log of the AuditLog.
	// example: {"log": {"request_logs": [{"id": "1", "status_logs": [{"status": "Pending", "timestamp": "2021-07-01T12:00
	// :00Z", "user": "John Doe", "user_id": "1"}], "order_logs": [{"id": "1", "status_logs": [{"status": "Pending", "timestamp": "2021-07-01T12:00
	// :00Z", "user": "John Doe", "user_id": "1"}], "sample_logs": [{"id": "1", "status_logs": [{"status": "Pending", "timestamp": "2021-07-01T12:00
	// :00Z", "user": "John Doe", "user_id": "1"}], "test_logs": [{"id": "1", "status_logs": [{"status": "Pending", "timestamp": "2021-07-01T12:00
	// :00Z", "user": "John Doe", "user_id": "1"}], "investigation_logs": [{"investigation_name": "CBC", "status_logs": [{"status": "Pending", "timestamp": "2021-07-01T12:00
	// :00Z", "user": "John Doe", "user_id": "1"}]}]}]}]}]}}
	Log Log `json:"log"`
}

// @swagger:response Log
type Log struct {
	// The request details of the Log.
	RequestLogs []RequestLog `json:"request_logs"`
}

// @swagger:response StatusLog
type StatusLog struct {
	// The status of the StatusLog.
	// example: "Pending"
	Status string `json:"status"`
	// The timestamp of the StatusLog.
	// example: "2021-07-01T12:00:00Z"
	Timestamp string `json:"timestamp"`
	// The user of the StatusLog.
	// example: "John Doe"
	User string `json:"user"`
	// The user ID of the StatusLog.
	// example: "1"
	UserId string `json:"user_id"`
}

// @swagger:response RequestLog
type RequestLog struct {
	// The ID of the RequestLog.
	// example: 1
	Id string `json:"id"`
	// The status logs of the RequestLog.
	// example: [{"status": "Pending", "timestamp": "2021-07-01T12:00
	// :00Z", "user": "John Doe", "user_id": "1"}]
	StatusLogs []StatusLog `json:"status_logs"`
	// The order logs of the RequestLog.
	// example: [{"id": "1", "status_logs": [{"status": "Pending", "timestamp": "2021-07-01T12:00
	// :00Z", "user": "John Doe", "user_id": "1"}], "sample_logs": [{"id": "1", "status_logs": [{"status": "Pending", "timestamp": "2021-07-01T12:00
	// :00Z", "user": "John Doe", "user_id": "1"}], "test_logs": [{"id": "1", "status_logs": [{"status": "Pending", "timestamp": "2021-07-01T12:00
	// :00Z", "user": "John Doe", "user_id": "1"}], "investigation_logs": [{"investigation_name": "CBC", "status_logs": [{"status": "Pending", "timestamp": "2021-07-01T12:00
	// :00Z", "user": "John Doe", "user_id": "1"}]}]}]}]
	OrderLogs []OrderLogs `json:"order_logs"`
}

// @swagger:response OrderLogs
type OrderLogs struct {
	// The ID of the OrderLogs.
	// example: 1
	Id string `json:"id"`
	// The status logs of the OrderLogs.
	// example: [{"status": "Pending", "timestamp": "2021-07-01T12:00
	// :00Z", "user": "John Doe", "user_id": "1"}]
	StatusLogs []StatusLog `json:"status_logs"`
	// The sample logs of the RequestLog.
	// example: [{"id": "1", "status_logs": [{"status": "Pending", "timestamp": "2021-07-01T12:00
	// :00Z", "user": "John Doe", "user_id": "1"}], "test_logs": [{"id": "1", "status_logs": [{"status": "Pending", "timestamp": "2021-07-01T12:00
	// :00Z", "user": "John Doe", "user_id": "1"}], "investigation_logs": [{"investigation_name": "CBC", "status_logs": [{"status": "Pending", "timestamp": "2021-07-01T12:00
	// :00Z", "user": "John Doe", "user_id": "1"}]}]}]
	SampleLogs []SampleLogs `json:"sample_logs"`
}

// @swagger:response SampleLogs
type SampleLogs struct {
	// The ID of the SampleLogs.
	// example: 1
	Id string `json:"id"`
	// The status logs of the SampleLogs.
	// example: [{"status": "Pending", "timestamp": "2021-07-01T12:00
	// :00Z", "user": "John Doe", "user_id": "1"}]
	StatusLogs []StatusLog `json:"status_logs"`
	// The test logs of the SampleLogs.
	// example: [{"id": "1", "status_logs": [{"status": "Pending", "timestamp": "2021-07-01T12:00
	// :00Z", "user": "John Doe", "user_id": "1"}], "investigation_logs": [{"investigation_name": "CBC", "status_logs": [{"status": "Pending", "timestamp": "2021-07-01T12:00
	// :00Z", "user": "John Doe", "user_id": "1"}]}]
	TestLogs []TestLogs `json:"test_logs"`
}

// @swagger:response TestLogs
type TestLogs struct {
	// The ID of the TestLogs.
	// example: 1
	Id string `json:"id"`
	// The status logs of the TestLogs.
	// example: [{"status": "Pending", "timestamp": "2021-07-01T12:00
	// :00Z", "user": "John Doe", "user_id": "1"}]
	StatusLogs []StatusLog `json:"status_logs"`
	// The investigation logs of the TestLogs.
	// example: [{"investigation_name": "CBC", "status_logs": [{"status": "Pending", "timestamp": "2021-07-01T12:00
	// :00Z", "user": "John Doe", "user_id": "1"}]
	InvestigationLogs []InvestigationLogs `json:"investigation_logs"`
}

// @swagger:response InvestigationLogs
type InvestigationLogs struct {
	// The investigation name of the InvestigationLogs.
	// example: "CBC"
	InvestigationName string `json:"investigation_name"`
	// The status logs of the InvestigationLogs.
	// example: [{"status": "Pending", "timestamp": "2021-07-01T12:00
	// :00Z", "user": "John Doe", "user_id": "1"}]
	StatusLogs []StatusLog `json:"status_logs"`
}

type SampleDBLog struct {
	// The operation applied - INSERT, UPDATE, DELETE
	Operation string `json:"operation"`
	// The column updated
	FieldName    string    `json:"field_name"`
	OldValue     string    `json:"old_value"`
	NewValue     string    `json:"new_value"`
	UserName     string    `json:"user_name"`
	LogTimestamp time.Time `json:"log_timestamp,omitempty"`
}

type SampleLogBody struct {
	SampleId   uint                        `json:"sample_id"`
	VialTypeId uint                        `json:"vial_type_id"`
	Logs       map[time.Time][]SampleDBLog `json:"logs"`
}
