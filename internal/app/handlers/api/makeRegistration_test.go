package api_handler

import (
	"bytes"
	"context"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"main/internal/app/firebase"
	fbMocks "main/internal/app/firebase/mocks"
	"main/internal/app/store/sqlstore"
	"main/internal/app/store/sqlstore/mocks"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewRegistrationHandler(t *testing.T) {
	type args struct {
		ctx   context.Context
		store sqlstore.StoreInterface
		log   *logrus.Logger
		fbAPI firebase.FirebaseAPIInterface
	}
	mockFirebase := fbMocks.NewFirebaseAPIInterface(t)
	mockStore := mocks.NewStoreInterface(t)
	apiMock := mocks.NewApiRepositoryInterface(t)
	mockStore.On("API").Return(apiMock)
	apiMock.On("RegistrationToken", mock.Anything, mock.Anything, mock.Anything).Return("", nil)
	argsStruct := args{
		ctx:   context.Background(),
		store: mockStore,
		log:   logrus.New(),
		fbAPI: mockFirebase,
	}
	//mockFirebase.On("GetFirebaseUser", mock.Anything, "test-case123").Return(&auth.UserRecord{}, nil)
	//mockFirebase.On("GetFirebaseUser", mock.Anything, "test-case_unauthorized").Return("HUI", errors.New("unauthorized"))
	handler := NewRegistrationHandler(argsStruct.ctx, argsStruct.store, argsStruct.log, argsStruct.fbAPI)
	goodRequest, _ := http.NewRequest("POST", "/api/token", bytes.NewBufferString("{\"device_tag\":\"test-case\",\"uid\":\"test-case123\"}"))
	badPayloadRequest, _ := http.NewRequest("POST", "/api/token", bytes.NewBufferString("{\"deviceTag\":\"test-case\",\"uid\":\"test-case123\"}"))
	badUIDRequest, _ := http.NewRequest("POST", "/api/token", bytes.NewBufferString("{\"deviceTag\":\"test-case\",\"uid\":\"case_unauthorized\"}"))
	tests := []struct {
		name               string
		wantStatus         int
		mockResponseWriter *httptest.ResponseRecorder
		mockRequest        *http.Request
		excpectCall        bool
	}{
		{
			name:               "ok-test",
			mockResponseWriter: httptest.NewRecorder(),
			mockRequest:        goodRequest,
			wantStatus:         http.StatusOK,
		},
		{
			name:               "bad-payload-test",
			mockResponseWriter: httptest.NewRecorder(),
			mockRequest:        badPayloadRequest,
			wantStatus:         http.StatusBadRequest,
		},
		{
			name:               "bad-uid-test",
			mockResponseWriter: httptest.NewRecorder(),
			mockRequest:        badUIDRequest,
			wantStatus:         http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler(tt.mockResponseWriter, tt.mockRequest)
			if gotStatus := tt.mockResponseWriter.Result().StatusCode; !reflect.DeepEqual(gotStatus, tt.wantStatus) {
				t.Errorf("NewRegistrationHandler() = %v, want status_code %v", gotStatus, tt.wantStatus)
			}
		})
	}
}
