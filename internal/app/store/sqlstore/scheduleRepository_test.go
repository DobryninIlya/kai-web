package sqlstore

import (
	"testing"
	"time"
)

func Test_isContainsInDict(t *testing.T) {
	type args struct {
		date string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "contains",
			args: args{date: "чет"},
			want: true,
		},
		{
			name: "contains",
			args: args{date: "чет/неч"},
			want: true,
		},
		{
			name: "not contains",
			args: args{date: "еженедельно"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isContainsInDict(tt.args.date); got != tt.want {
				t.Errorf("isContainsInDict() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isContainDate(t *testing.T) {
	type args struct {
		data   string
		margin int
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		// TODO: Add test cases.
		{
			name: "contains",
			args: args{
				data:   time.Now().AddDate(0, 0, 0).Format("02.01"),
				margin: 0,
			},
			want:  time.Now().AddDate(0, 0, 0).Format("02.01"),
			want1: time.Now().AddDate(0, 0, 0).Format("2.01"),
		},
		{
			name: "do not contains",
			args: args{
				data:   time.Now().AddDate(0, 0, 0).Format("02.01"),
				margin: 1,
			},
			want:  "",
			want1: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := isContainDate(tt.args.data, tt.args.margin)
			if got != tt.want || got1 != tt.want1 {
				t.Errorf("isContainDate() got = %v, want %v", got, tt.want)
			}

		})
	}
}

func Test_getSubgroupForDate(t *testing.T) {
	type args struct {
		data   string
		ex1    string
		ex2    string
		isEven bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				data:   "2.06 / 20.06",
				ex1:    "02.06",
				ex2:    "2.06",
				isEven: true,
			},
			want: "[1 гр.]",
		},
		{
			name: "",
			args: args{
				data:   "2.06 / 20.06",
				ex1:    "02.06",
				ex2:    "2.06",
				isEven: false,
			},
			want: "[2 гр.]",
		},
		{
			name: "",
			args: args{
				data:   "2.06, 04.06 / 20.06, 25.06",
				ex1:    "4.06",
				ex2:    "04.06",
				isEven: true,
			},
			want: "[1 гр.]",
		},
		{
			name: "",
			args: args{
				data:   "2.06, 04.06 / 20.06, 25.06",
				ex1:    "20.06",
				ex2:    "20.06",
				isEven: false,
			},
			want: "[1 гр.]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSubgroupForDate(tt.args.data, tt.args.ex1, tt.args.ex2, tt.args.isEven); got != tt.want {
				t.Errorf("getSubgroupForDate() = %v, want %v", got, tt.want)
			}
		})
	}
}
