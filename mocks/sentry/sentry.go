package sentryMocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockSentryDependency struct {
	mock.Mock
}

func (m *MockSentryDependency) LogError(ctx context.Context, message string, err error, additionalInformation map[string]interface{}) {
	m.Called(message, err, additionalInformation)
}
