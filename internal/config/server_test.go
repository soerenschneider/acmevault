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
					Addr:         "https://vault:8200",
					SecretId:     "secretId",
					RoleId:       "roleId",
					PathPrefix:   "preprod",
					AuthMethod:   "approle",
					Kv2MountPath: "secret",
					AwsMountPath: "custom-aws-mountpath",
					AwsRole:      "my-custom-role",
				},
				AcmeEmail:       "my@email.tld",
				AcmeUrl:         letsEncryptUrl,
				AcmeDnsProvider: "",
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

func TestAcmeVaultServerConfig_Validate(t *testing.T) {
	type fields struct {
		VaultConfig          VaultConfig
		AcmeEmail            string
		AcmeUrl              string
		AcmeDnsProvider      string
		AcmeCustomDnsServers []string
		IntervalSeconds      int
		Domains              []AcmeServerDomains
		MetricsAddr          string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid example",
			fields: fields{
				VaultConfig: VaultConfig{
					Token:            "token",
					Addr:             "https://my-vault",
					PathPrefix:       "bla",
					DomainPathFormat: "blub-%s",
					AuthMethod:       "approle",
					Kv2MountPath:     "secret",
					AwsMountPath:     "custom-aws-mountpath",
					AwsRole:          "my-custom-role",
					SecretId:         "secret-id",
					RoleId:           "role",
				},
				AcmeEmail:            "ac@me.com",
				AcmeUrl:              letsEncryptUrl,
				AcmeDnsProvider:      "",
				AcmeCustomDnsServers: []string{"8.8.8.8", "2001:4860:4860::8888"},
				IntervalSeconds:      3600,
				Domains: []AcmeServerDomains{
					{
						Domain: "valid.domain",
						Sans:   []string{"another.valid.domain"},
					},
				},
				MetricsAddr: "127.0.0.1:9100",
			},
			wantErr: false,
		},
		{
			name: "invalid custom dns servers",
			fields: fields{
				VaultConfig: VaultConfig{
					Token:            "token",
					Addr:             "https://my-vault",
					PathPrefix:       "bla",
					DomainPathFormat: "blub-%s",
					AuthMethod:       "approle",
					Kv2MountPath:     "secret",
				},
				AcmeEmail:            "ac@me.com",
				AcmeUrl:              letsEncryptUrl,
				AcmeDnsProvider:      "",
				AcmeCustomDnsServers: []string{"not.an.ip"},
				IntervalSeconds:      3600,
				Domains: []AcmeServerDomains{
					{
						Domain: "valid.domain",
						Sans:   []string{"another.valid.domain"},
					},
				},
				MetricsAddr: "127.0.0.1:9100",
			},
			wantErr: true,
		},
		{
			name: "domain has no valid fqdn",
			fields: fields{
				VaultConfig: VaultConfig{
					Token:            "token",
					Addr:             "https://my-vault",
					PathPrefix:       "bla",
					DomainPathFormat: "blub-%s",
					AuthMethod:       "token",
					Kv2MountPath:     "secret",
				},
				AcmeEmail:            "ac@me.com",
				AcmeUrl:              letsEncryptUrl,
				AcmeDnsProvider:      "",
				AcmeCustomDnsServers: nil,
				IntervalSeconds:      3600,
				Domains: []AcmeServerDomains{
					{
						Domain: "nofqdn",
						Sans:   []string{"another.valid.domain"},
					},
				},
				MetricsAddr: "127.0.0.1:9100",
			},
			wantErr: true,
		},
		{
			name: "sans has no valid fqdn",
			fields: fields{
				VaultConfig: VaultConfig{
					Token:            "token",
					Addr:             "https://my-vault",
					PathPrefix:       "bla",
					DomainPathFormat: "blub-%s",
					AuthMethod:       "token",
					Kv2MountPath:     "secret",
				},
				AcmeEmail:            "ac@me.com",
				AcmeUrl:              letsEncryptUrl,
				AcmeDnsProvider:      "",
				AcmeCustomDnsServers: nil,
				IntervalSeconds:      3600,
				Domains: []AcmeServerDomains{
					{
						Domain: "valid.fqdn",
						Sans:   []string{"novalidfqdn", "valid.fqdn"},
					},
				},
				MetricsAddr: "127.0.0.1:9100",
			},
			wantErr: true,
		},
		{
			name: "interval seconds too low",
			fields: fields{
				VaultConfig: VaultConfig{
					Token:            "token",
					Addr:             "https://my-vault",
					PathPrefix:       "bla",
					DomainPathFormat: "blub-%s",
					AuthMethod:       "token",
					Kv2MountPath:     "secret",
				},
				AcmeEmail:            "ac@me.com",
				AcmeUrl:              letsEncryptUrl,
				AcmeDnsProvider:      "",
				AcmeCustomDnsServers: nil,
				IntervalSeconds:      3599,
				Domains: []AcmeServerDomains{
					{
						Domain: "valid.fqdn",
						Sans:   []string{"one.more.valid.fqdn", "another.valid.fqdn"},
					},
				},
				MetricsAddr: "127.0.0.1:9100",
			},
			wantErr: true,
		},
		{
			name: "interval seconds too high",
			fields: fields{
				VaultConfig: VaultConfig{
					Token:            "token",
					Addr:             "https://my-vault",
					PathPrefix:       "bla",
					DomainPathFormat: "blub-%s",
					AuthMethod:       "token",
					Kv2MountPath:     "secret",
				},
				AcmeEmail:            "ac@me.com",
				AcmeUrl:              letsEncryptUrl,
				AcmeDnsProvider:      "",
				AcmeCustomDnsServers: nil,
				IntervalSeconds:      86401,
				Domains: []AcmeServerDomains{
					{
						Domain: "valid.fqdn",
						Sans:   []string{"one.more.valid.fqdn", "another.valid.fqdn"},
					},
				},
				MetricsAddr: "127.0.0.1:9100",
			},
			wantErr: true,
		},
		{
			name: "invalid acme email",
			fields: fields{
				VaultConfig: VaultConfig{
					Token:            "token",
					Addr:             "https://my-vault",
					PathPrefix:       "bla",
					DomainPathFormat: "blub-%s",
					AuthMethod:       "token",
					Kv2MountPath:     "secret",
					AwsMountPath:     "custom-aws-mountpath",
					AwsRole:          "my-custom-role",
				},
				AcmeEmail:            "bla",
				AcmeUrl:              letsEncryptUrl,
				AcmeDnsProvider:      "",
				AcmeCustomDnsServers: nil,
				IntervalSeconds:      86400,
				Domains: []AcmeServerDomains{
					{
						Domain: "valid.fqdn",
						Sans:   []string{"one.more.valid.fqdn", "another.valid.fqdn"},
					},
				},
				MetricsAddr: "127.0.0.1:9100",
			},
			wantErr: true,
		},
		{
			name: "invalid acme url",
			fields: fields{
				VaultConfig: VaultConfig{
					Token:            "token",
					Addr:             "https://my-vault",
					PathPrefix:       "bla",
					DomainPathFormat: "blub-%s",
					AuthMethod:       "token",
					Kv2MountPath:     "secret",
				},
				AcmeEmail:            "ac@me.com",
				AcmeUrl:              "not valid url!",
				AcmeDnsProvider:      "",
				AcmeCustomDnsServers: nil,
				IntervalSeconds:      86400,
				Domains: []AcmeServerDomains{
					{
						Domain: "valid.fqdn",
						Sans:   []string{"one.more.valid.fqdn", "another.valid.fqdn"},
					},
				},
				MetricsAddr: "127.0.0.1:9100",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := AcmeVaultServerConfig{
				VaultConfig:          tt.fields.VaultConfig,
				AcmeEmail:            tt.fields.AcmeEmail,
				AcmeUrl:              tt.fields.AcmeUrl,
				AcmeDnsProvider:      tt.fields.AcmeDnsProvider,
				AcmeCustomDnsServers: tt.fields.AcmeCustomDnsServers,
				IntervalSeconds:      tt.fields.IntervalSeconds,
				Domains:              tt.fields.Domains,
				MetricsAddr:          tt.fields.MetricsAddr,
			}
			if err := conf.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
