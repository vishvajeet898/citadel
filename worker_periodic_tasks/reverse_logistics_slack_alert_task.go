package workerPeriodicTasks

import (
	"context"
	"fmt"
	"time"

	"github.com/RichardKnop/machinery/v1/tasks"

	"github.com/Orange-Health/citadel/common/constants"
)

func (s *WorkerPeriodicTaskService) SlackAlertForDelayedReverseLogisticsPeriodicTask() error {
	ctx := context.Background()
	s.SampleService.SendSlackAlertForDelayedReverseLogisticsTat(ctx)
	s.SampleService.SendSlackAlertForDelayedInterlabLogisticsTat(ctx)
	return nil
}

func SlackAlertForDelayedReverseLogisticsPeriodicTaskSignature() *tasks.Signature {
	groupId := fmt.Sprintf("%s:%v", constants.SlackAlertForReverseLogisticsPeriodicTaskName, time.Now().Unix())
	return &tasks.Signature{
		Name:                 constants.SlackAlertForReverseLogisticsPeriodicTaskName,
		Args:                 nil,
		RoutingKey:           constants.WorkerDefaultQueue,
		BrokerMessageGroupId: groupId,
	}
}
