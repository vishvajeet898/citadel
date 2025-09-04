package workerPeriodicTasks

import (
	"context"
	"fmt"
	"time"

	"github.com/RichardKnop/machinery/v1/tasks"

	"github.com/Orange-Health/citadel/common/constants"
)

func (s *WorkerPeriodicTaskService) EtsTatBreachedPeriodicTask() error {
	ctx := context.Background()
	s.EtsService.HandleTatBreachedTestsCronByEvents(ctx)

	return nil
}

func EtsTatBreachedPeriodicTaskSignature() *tasks.Signature {
	groupId := fmt.Sprintf("%s:%v", constants.EtsTatBreachedPeriodicTaskName, time.Now().Unix())
	return &tasks.Signature{
		Name:                 constants.EtsTatBreachedPeriodicTaskName,
		Args:                 nil,
		RoutingKey:           constants.WorkerDefaultQueue,
		BrokerMessageGroupId: groupId,
	}
}
