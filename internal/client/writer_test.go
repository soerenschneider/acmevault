package client

import (
	"runtime"
	"testing"
)

func Test_getUidFromUsername(t *testing.T) {
	type args struct {
		username string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "existent",
			args:    args{"root"},
			want:    0,
			wantErr: false,
		},
		{
			name:    "nonexistent",
			args:    args{"ihopeyoudontexist"},
			want:    -1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getUidFromUsername(tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("getUidFromUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getUidFromUsername() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getGidFromGroup(t *testing.T) {
	type args struct {
		group string
	}
	tests := []struct {
		os      string
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			os:      "darwin",
			name:    "existent-darwin",
			args:    args{"wheel"},
			want:    0,
			wantErr: false,
		},
		{
			os:      "linux",
			name:    "existent-linux",
			args:    args{"root"},
			want:    0,
			wantErr: false,
		},
		{
			name:    "nonexistent",
			args:    args{"ihopeyoudontexist"},
			want:    -1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		if runtime.GOOS == tt.os || tt.os == "" {
			t.Run(tt.name, func(t *testing.T) {
				got, err := getGidFromGroup(tt.args.group)
				if (err != nil) != tt.wantErr {
					t.Errorf("getGidFromGroup() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("getGidFromGroup() got = %v, want %v", got, tt.want)
				}
			})
		}
	}
}

func Test_runHooks(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "nil",
			args:    nil,
			wantErr: false,
		},
		{
			name:    "empty",
			args:    []string{},
			wantErr: false,
		},
		{
			name:    "cmd",
			args:    []string{"date"},
			wantErr: false,
		},
		{
			name:    "cmd with arg",
			args:    []string{"date", "+%s"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := executeHook(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("executeHook() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMd5Compare(t *testing.T) {
	type args struct {
		a []byte
		b []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty",
			args: args{
				a: nil,
				b: nil,
			},
			want: true,
		},
		{
			name: "one empty",
			args: args{
				a: []byte("hello"),
				b: nil,
			},
			want: false,
		},
		{
			name: "equal",
			args: args{
				a: []byte("hello"),
				b: []byte("hello"),
			},
			want: true,
		},
		{
			name: "not equal",
			args: args{
				a: []byte("hello"),
				b: []byte("world"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Md5Compare(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Md5Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}
