package config

import (
	"reflect"
	"testing"
)

func TestAcmeVaultServerConfigFromFile(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    AcmeVaultServerConfig
		wantErr bool
	}{
		{
			name: "example server config",
			path: "../../contrib/server.json",
			want: AcmeVaultServerConfig{
				VaultConfig: VaultConfig{
					VaultAddr:  "https://vault:8200",
					SecretId:   "secretId",
					RoleId:     "roleId",
					PathPrefix: "preprod",
				},
				AcmeConfig: AcmeConfig{
					Email:           "my@email.tld",
					AcmeUrl:         letsEncryptUrl,
					AcmeDnsProvider: "",
				},
				IntervalSeconds: 43200,
				Domains: []AcmeServerDomains{
					{
						Domain: "domain1.tld",
						Sans:   []string{"domain3.tld", "domain4.tld"},
					},
					{
						Domain: "domain2.tld",
					},
				},
				MetricsAddr: "127.0.0.1:9112",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AcmeVaultServerConfigFromFile(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("AcmeVaultServerConfigFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AcmeVaultServerConfigFromFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}
