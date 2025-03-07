package handlers_test

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kiricle/api-homework/internal/handlers"
	mock_handlers "github.com/kiricle/api-homework/internal/handlers/mocks"
	"github.com/kiricle/api-homework/internal/models"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBookHandler_GetBooks(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookService := mock_handlers.NewMockBookService(ctrl)

	mockBooks := []models.Book{
		{
			Id:     1,
			Name:   "Norway wood",
			Author: "Murakami",
		},
		{
			Id:     2,
			Name:   "Nothing happens",
			Author: "Kyrylo",
		},
	}

	mockBookService.EXPECT().GetBooks().Return(mockBooks, nil)

	logger := slog.Default()
	bookHandler := handlers.NewBookHandler(logger, mockBookService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	bookHandler.GetBooks(c)

	assert.Equal(t, w.Code, http.StatusOK)
	assert.Contains(t, w.Body.String(), "Norway wood")
	assert.Contains(t, w.Body.String(), "Nothing happens")
}
