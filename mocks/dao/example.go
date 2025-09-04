package daoMocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockExampleDaoDependency struct {
	mock.Mock
}

func (m *MockExampleDaoDependency) Example(ctx context.Context) string {
	args := m.Called()
	return args.String(0)
}
