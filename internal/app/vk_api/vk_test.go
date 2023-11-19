package vk_api

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"testing"
)

type RoundTripFunc func(req *http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestAPIvk_SendMessageVKids(t *testing.T) {
	type fields struct {
		vkToken    string
		vkTemplate string
		httpClient http.Client
	}
	type args struct {
		log     *logrus.Logger
		uId     []int64
		message string
		buttons string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "test1",
			fields: fields{
				vkToken:    "VK_TEST_TOKEN",
				vkTemplate: vkTemplate,
				httpClient: *&http.Client{
					Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
						response := http.Response{
							StatusCode: http.StatusOK,
						}
						return &response, nil
					}),
				},
			},
			args: args{
				log:     logrus.New(),
				uId:     []int64{1, 2, 3},
				message: "test",
				buttons: "test",
			},
			want: true,
		},
		{
			name: "test2",
			fields: fields{
				vkToken:    "VK_TEST_FALL_TOKEN",
				vkTemplate: vkTemplate,
				httpClient: *&http.Client{
					Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
						response := http.Response{
							StatusCode: http.StatusForbidden,
						}
						return &response, nil
					}),
				},
			},
			args: args{
				log:     logrus.New(),
				uId:     []int64{1, 2, 3},
				message: "test",
				buttons: "test",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := APIvk{
				vkToken:    tt.fields.vkToken,
				vkTemplate: tt.fields.vkTemplate,
				httpClient: tt.fields.httpClient,
			}
			if got := r.SendMessageVKids(tt.args.log, tt.args.uId, tt.args.message, tt.args.buttons); got != tt.want {
				t.Errorf("SendMessageVKids() = %v, want %v", got, tt.want)
			}
		})
	}
}
