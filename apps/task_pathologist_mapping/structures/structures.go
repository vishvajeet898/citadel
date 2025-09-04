package structures

// @swagger: response taskPathologistMapping
type TaskPathologistMapping struct {
	// The ID of the task.
	// example: 1
	TaskID uint `json:"task_id"`
	// The ID of the pathologist.
	// example: 1
	PathologistID uint `json:"pathologist_id"`
	// Whether the task pathologist mapping is active.
	// example: true
	IsActive bool `json:"is_active"`
}
