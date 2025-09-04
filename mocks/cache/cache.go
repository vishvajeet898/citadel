package cacheMocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockCacheDependency struct {
	mock.Mock
}

func (m *MockCacheDependency) Set(ctx context.Context, key string, value interface{}, expiration int) error {
	m.Called(ctx, key, value, expiration)
	return nil
}

func (m *MockCacheDependency) Get(ctx context.Context, key string, result interface{}) error {
	m.Called(ctx, key, result)
	return nil
}

func (m *MockCacheDependency) Exists(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func (m *MockCacheDependency) Delete(ctx context.Context, key string) error {
	m.Called(ctx, key)
	return nil
}

func (m *MockCacheDependency) HSet(ctx context.Context, key, field string, value interface{}, expiry time.Duration) error {
	m.Called(ctx, key, field, value, expiry)
	return nil
}

func (m *MockCacheDependency) HSetAll(ctx context.Context, key string, value interface{}, expiry time.Duration) error {
	m.Called(ctx, key, value, expiry)
	return nil
}

func (m *MockCacheDependency) HGet(ctx context.Context, key, field string) (string, error) {
	args := m.Called(ctx, key, field)
	return args.String(0), args.Error(1)
}

func (m *MockCacheDependency) HGetAll(ctx context.Context, key string) (map[string]interface{}, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}
