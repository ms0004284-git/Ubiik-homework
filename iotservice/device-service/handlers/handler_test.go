package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"errors"

	"device-service/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetDeviceByID(id string) (*models.Device, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Device), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) CreateDevice(device *models.Device) error {
	args := m.Called(device)
	return args.Error(0)
}

func (m *MockRepository) UpdateDevice(device *models.Device) error {
	args := m.Called(device)
	return args.Error(0)
}

func TestUpdateOrCreateDevice(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockRepo := new(MockRepository)
	handler := NewDeviceAccess(mockRepo)

	r.PUT("/devices/:deviceId", handler.UpdateOrCreateDevice)
	t.Run("Create new device", func(t *testing.T) {
		mockRepo.On("GetDeviceByID", "123").Return(nil, errors.New("not found"))
		mockRepo.On("CreateDevice", mock.Anything).Return(nil)

		body := `{"username":"testuser"}`
		req := httptest.NewRequest(http.MethodPut, "/devices/123", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		t.Logf("Response Body: %s", w.Body.String())

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), `"message":"Device created"`)

		mockRepo.AssertExpectations(t)
		mockRepo.ExpectedCalls = nil // 清除 mockRepo 的暫存
	})

	t.Run("Update existing device", func(t *testing.T) {
		existingDevice := &models.Device{ID: "123", Username: "olduser"}
		mockRepo.On("GetDeviceByID", "123").Return(existingDevice, nil)
		mockRepo.On("UpdateDevice", mock.Anything).Return(nil)

		body := `{"username":"updateduser"}`
		req := httptest.NewRequest(http.MethodPut, "/devices/123", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"message":"Device username updated."`)

		mockRepo.AssertExpectations(t)

	})
}

func TestGetDeviceByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockRepo := new(MockRepository)
	handler := NewDeviceAccess(mockRepo)

	r.GET("/devices/:deviceId/username", handler.GetDevice)
	t.Run("Device exists", func(t *testing.T) {
		existingDevice := &models.Device{ID: "123", Username: "user"}
		mockRepo.On("GetDeviceByID", "123").Return(existingDevice, nil)

		req := httptest.NewRequest(http.MethodGet, "/devices/123/username", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"username":"user"`)

		mockRepo.AssertExpectations(t)
		mockRepo.ExpectedCalls = nil
	})

	t.Run("Device not found", func(t *testing.T) {
		mockRepo.On("GetDeviceByID", "123").Return(nil, errors.New("not found"))

		req := httptest.NewRequest(http.MethodGet, "/devices/123/username", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), `"message":"Item not found."`)

		mockRepo.AssertExpectations(t)
	})
}