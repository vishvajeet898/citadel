package tests

// import (
// 	"context"
// 	"testing"

// 	"github.com/Orange-Health/citadel/apps/example/service"
// 	cacheMocks "github.com/Orange-Health/citadel/mocks/cache"
// 	daoMocks "github.com/Orange-Health/citadel/mocks/dao"
// 	sentryMocks "github.com/Orange-Health/citadel/mocks/sentry"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// func TestExampleFunction(t *testing.T) {
// 	dbDependency := new(daoMocks.MockExampleDaoDependency)
// 	dbDependency.On("Example").Return("mocked value")

// 	cacheMock := new(cacheMocks.MockCacheDependency)
// 	cacheMock.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
// 	cacheMock.On("Get", mock.Anything, mock.Anything, mock.Anything)

// 	sentryMock := new(sentryMocks.MockSentryDependency)
// 	sentryMock.On("LogError", mock.Anything, mock.Anything, mock.Anything)

// 	service := &service.ExampleService{
// 		Db:     dbDependency,
// 		Cache:  cacheMock,
// 		Sentry: sentryMock,
// 	}

// 	ctx := context.Background()
// 	result := service.Example(ctx)

// 	assert.Equal(t, "Example mocked value", result)
// }
