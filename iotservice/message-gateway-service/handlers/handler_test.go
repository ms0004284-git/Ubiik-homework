package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"message-gateway-service/models"
	"github.com/stretchr/testify/assert"
)

func TestHandleMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		message        interface{}
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Register Device",
			message: models.Message{
				DeviceID: "my-device",
				Type:     "registration",
				Data: models.RegistrationData{
					Username: "username",
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
		{
			name: "Unknown message type",
			message: models.Message{
				DeviceID: "12345",
				Type:     "unknown",
				Data: models.RegistrationData{
					Username: "username",
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Unknown message type."}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.message)
			req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
			ctx.Request = req

			HandleMessage(ctx)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}