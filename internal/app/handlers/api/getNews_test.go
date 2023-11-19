package api_handler

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"main/internal/app/model"
	"main/internal/app/store/sqlstore/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewNewsHandler(t *testing.T) {
	mockStore := mocks.NewStoreInterface(t)
	apiMock := mocks.NewApiRepositoryInterface(t)
	mockLogger := logrus.New()
	mockResponseWriter := httptest.NewRecorder()
	mockRequest, _ := http.NewRequest("GET", "/api/news/20", nil)
	mockRequest = mockRequest.WithContext(context.WithValue(mockRequest.Context(), chi.RouteCtxKey, &chi.Context{
		URLParams: chi.RouteParams{
			Keys:   []string{"newsId"},
			Values: []string{"123"},
		},
	}))
	// Создаем экземпляр хэндлера
	handler := NewNewsHandler(mockStore, mockLogger)

	// Устанавливаем ожидаемые вызовы и возвращаемые значения для mockStore
	//expectedNewsID := 123
	mockStore.On("API").Return(apiMock)
	apiMock.On("GetNewsById", 123).Return(model.News{}, nil)
	// Вызываем хэндлер
	handler(mockResponseWriter, mockRequest)

	// Проверяем ожидаемый статус код
	//expectedStatusCode := http.StatusOK
	//if mockResponseWriter.Result().StatusCode != 0 { // != expectedStatusCode {
	//	t.Errorf("Expected status code %d, got %d", expectedStatusCode, mockResponseWriter.Result().StatusCode)
	//}

	// Проверяем ожидаемое содержимое ответа

	//mockStore.AssertCalled(t, "API")
	//mockStore.AssertCalled(t, "GetNewsById", 123)
}
