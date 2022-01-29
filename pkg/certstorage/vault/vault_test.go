package vault

import (
	"github.com/hashicorp/vault/api"
	"github.com/soerenschneider/acmevault/internal/config"
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

func TestVaultBackend_getSecretDataPath(t *testing.T) {
	type fields struct {
		client           *api.Client
		conf             config.VaultConfig
		revokeToken      bool
		namespacedPrefix string
	}
	tests := []struct {
		name   string
		fields fields
		args   string
		want   string
	}{
		{
			name: "custom domain path",
			fields: fields{
				client: nil,
				conf: config.VaultConfig{
					DomainPathFormat: "machine-%s",
				},
				revokeToken:      false,
				namespacedPrefix: "acmevault",
			},
			args: "test.domain.tld",
			want: "acmevault/client/machine-test.domain.tld/privatekey",
		},
		{
			name: "no domain format given",
			fields: fields{
				client:           nil,
				conf:             config.VaultConfig{},
				revokeToken:      false,
				namespacedPrefix: "acmevault",
			},
			args: "test.domain.tld",
			want: "acmevault/client/test.domain.tld/privatekey",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vault := &VaultBackend{
				client:           tt.fields.client,
				conf:             tt.fields.conf,
				revokeToken:      tt.fields.revokeToken,
				namespacedPrefix: tt.fields.namespacedPrefix,
			}
			if got := vault.getSecretDataPath(tt.args); got != tt.want {
				t.Errorf("getSecretDataPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
