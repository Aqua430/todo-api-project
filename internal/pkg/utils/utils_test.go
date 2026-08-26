package utils_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"todo-api/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

func TestMustGetID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		paramValue     string
		expectedID     int
		expectedOk     bool
		expectedStatus int
	}{
		{
			name:           "Valid positive ID",
			paramValue:     "42",
			expectedID:     42,
			expectedOk:     true,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Non-numeric ID",
			paramValue:     "abc",
			expectedID:     0,
			expectedOk:     false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Zero ID",
			paramValue:     "0",
			expectedID:     0,
			expectedOk:     false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Negative ID",
			paramValue:     "-10",
			expectedID:     0,
			expectedOk:     false,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Params = []gin.Param{{Key: "id", Value: tt.paramValue}}

			id, ok := utils.MustGetID(c, "id")

			if ok != tt.expectedOk {
				t.Errorf("got ok = %v, want %v", ok, tt.expectedOk)
			}
			if id != tt.expectedID {
				t.Errorf("got id = %d, want %d", id, tt.expectedID)
			}
			if w.Code != tt.expectedStatus {
				t.Errorf("got status = %d, want %d", w.Code, tt.expectedStatus)
			}
		})
	}
}
