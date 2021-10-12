package vault

import (
	"reflect"
	"testing"
)

func Test_buildSecretPayload(t *testing.T) {
	type args struct {
		data map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "simple",
			args: args{
				data: map[string]interface{}{
					"test": "bla",
				},
			},
			want:    []byte("{\"data\":{\"test\":\"bla\"},\"options\":{\"max_versions\":1}}"),
			wantErr: false,
		},
		{
			name:    "empty",
			args:    args{},
			want:    []byte("{\"data\":{},\"options\":{\"max_versions\":1}}"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := wrapPayload(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("wrapPayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("wrapPayload() got = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
