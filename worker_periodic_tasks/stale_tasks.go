package workerPeriodicTasks

import (
	"context"
	"fmt"
	"time"

	"github.com/RichardKnop/machinery/v1/tasks"
	"gorm.io/gorm"

	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/common/structures"
	"github.com/Orange-Health/citadel/common/utils"
)

func (s *WorkerPeriodicTaskService) StaleTasksPeriodicTask() error {
	ctx := context.Background()

	taskIdPreviousStateMap := []structures.TaskIdPreviousStateStruct{}
	s.Db.WithContext(ctx).Table(constants.TableTasks).
		Joins("INNER JOIN task_metadata ON tasks.id = task_metadata.task_id").
		Where("tasks.status = ?", constants.TASK_STATUS_IN_PROGRESS).
		Where("task_metadata.last_event_sent_at < ?", time.Now().Add(time.Duration(-1*constants.StaleTaskDuration)*time.Minute)).
		Select("tasks.id as task_id, tasks.previous_status").
		Find(&taskIdPreviousStateMap)

	currentTime := time.Now()
	taskUpdates := map[string]interface{}{
		"updated_at": &currentTime,
	}

	taskPathologistUpdates := map[string]interface{}{
		"updated_at": &currentTime,
		"is_active":  false,
	}
	err := s.Db.Transaction(func(tx *gorm.DB) error {
		for _, taskStruct := range taskIdPreviousStateMap {
			if taskStruct.PreviousStatus != "" {
				taskUpdates["status"] = taskStruct.PreviousStatus
				err := s.Db.WithContext(ctx).Table(constants.TableTasks).
					Where("id = ?", taskStruct.TaskId).
					Updates(taskUpdates).Error
				if err != nil {
					return err
				}

				err = s.Db.WithContext(ctx).Table(constants.TableTaskPathologistMapping).
					Where("task_id = ?", taskStruct.TaskId).
					Updates(taskPathologistUpdates).Error
				if err != nil {
					return err
				}
			}
		}
		return nil
	})

	if err != nil {
		utils.AddLog(ctx, constants.ERROR_LEVEL, utils.GetCurrentFunctionName(), nil, err)
		return err
	}

	return nil
}

func StaleTasksPeriodicTaskSignature() *tasks.Signature {
	groupId := fmt.Sprintf("%s:%v", constants.StaleTasksPeriodicTask, time.Now().Unix())
	return &tasks.Signature{
		Name:                 constants.StaleTasksPeriodicTask,
		Args:                 nil,
		RoutingKey:           constants.WorkerDefaultQueue,
		BrokerMessageGroupId: groupId,
	}
}
