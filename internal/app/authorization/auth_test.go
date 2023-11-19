package authorization

import (
	"testing"
)

func TestAuthorization_GetCookiesByPassword(t *testing.T) {
	type args struct {
		login    string
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr error
		wantLen int
	}{
		{
			name: "testFail",
			args: args{
				login:    "test",
				password: "test",
			},
			want:    false,
			wantErr: ErrWrongPassword,
			wantLen: 0,
		},
		{
			name: "testTrue",
			args: args{
				login:    "DobryninIS",
				password: "a4c13shj",
			},
			want:    false,
			wantErr: nil,
			wantLen: 7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Authorization{
				Cookies: *NewSafeMap(),
			}
			cookies, err := r.GetCookiesByPassword(tt.args.login, tt.args.password)
			if err != tt.wantErr {
				t.Errorf("CheckPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(cookies) != tt.wantLen {
				t.Errorf("CheckPassword() len(cookies) got = %v, want %v", len(cookies), 7)
			}
		})
	}
}
