package vk_app

//
//import (
//	"main/internal/app/handlers"
//	"main/internal/app/store/sqlstore"
//	"reflect"
//	"testing"
//)
//
//type storeMock struct {
//	sqlstore.StoreInterface
//}
//
//func TestNewRegistrationHandler(t *testing.T) {
//	store := storeMock{}
//	s := newApp(store)
//	tests := []struct {
//		name         string
//		payload      interface{}
//		expectedCode int
//	}{
//		// TODO: Add test cases.
//		{
//			name: "valid",
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			//s.ServeHTTP(reflect.interface{})
//			if got := handler.NewRegistrationHandler(storeMock{}); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("NewRegistrationHandler() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
