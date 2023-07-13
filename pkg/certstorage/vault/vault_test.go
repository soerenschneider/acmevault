package vault

import (
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/soerenschneider/acmevault/internal/config"
)

func TestVaultBackend_getSecretDataPath(t *testing.T) {
	type fields struct {
		client           *api.Client
		conf             config.VaultConfig
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
				namespacedPrefix: "acmevault",
			},
			args: "test.domain.tld",
			want: "acmevault/client/test.domain.tld/privatekey",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vault := &VaultBackend{
				client:   tt.fields.client,
				conf:     tt.fields.conf,
				basePath: tt.fields.namespacedPrefix,
			}
			if got := vault.getSecretDataPath(tt.args); got != tt.want {
				t.Errorf("getSecretDataPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVaultBackend_writeKv2Secret(t *testing.T) {

}
