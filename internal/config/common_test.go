package config

import "testing"

func TestVaultConfig_Validate(t *testing.T) {
	type fields struct {
		VaultToken       string
		VaultAddr        string
		SecretId         string
		RoleId           string
		PathPrefix       string
		SecretIdFile     string
		DomainPathFormat string
		AuthMethod       string
		Kv2MountPath     string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid config - token",
			fields: fields{
				AuthMethod:   "token",
				VaultToken:   "s.asd83hrfhasfjsda",
				VaultAddr:    "https://my-vault-instance:443",
				PathPrefix:   "production",
				Kv2MountPath: "secret",
			},
		},
		{
			name: "valid config - approle",
			fields: fields{
				AuthMethod:   "approle",
				VaultAddr:    "https://my-vault-instance:443",
				SecretId:     "super-secret",
				RoleId:       "my-role",
				PathPrefix:   "dev-v002",
				Kv2MountPath: "secret",
			},
		},
		{
			name: "valid config - approle, secret_id file",
			fields: fields{
				AuthMethod:   "approle",
				VaultAddr:    "https://my-vault-instance:443",
				SecretIdFile: "super-secret",
				RoleId:       "my-role",
				PathPrefix:   "dev-v002",
				Kv2MountPath: "secret",
			},
		},
		{
			name: "invalid config - missing protocol",
			fields: fields{
				AuthMethod:   "token",
				VaultToken:   "s.asd83hrfhasfjsda",
				VaultAddr:    "my-vault-instance:443",
				PathPrefix:   "production",
				Kv2MountPath: "secret",
			},
			wantErr: true,
		},
		{
			name: "invalid config - invalid path prefix",
			fields: fields{
				AuthMethod: "token",

				VaultToken:   "s.asd83hrfhasfjsda",
				VaultAddr:    "http://my-vault-instance:443",
				PathPrefix:   "/production",
				Kv2MountPath: "secret",
			},
			wantErr: true,
		},
		{
			name: "invalid config - no auth methods",
			fields: fields{
				AuthMethod: "approle",

				VaultAddr:    "http://my-vault-instance:443",
				PathPrefix:   "production",
				Kv2MountPath: "secret",
			},
			wantErr: true,
		},
		{
			name: "invalid config - empty path prefix",
			fields: fields{
				AuthMethod: "token",

				VaultAddr:    "http://my-vault-instance:443",
				VaultToken:   "s.VALIDVALIDVALID",
				PathPrefix:   "",
				Kv2MountPath: "secret",
			},
			wantErr: true,
		},
		{
			name: "invalid config - specifying secretId and secretIdFile",
			fields: fields{
				AuthMethod: "approle",

				VaultAddr:    "http://my-vault-instance:443",
				VaultToken:   "s.VALIDVALIDVALID",
				PathPrefix:   "production",
				RoleId:       "role",
				SecretId:     "secret-id",
				SecretIdFile: "/tmp/secret-id",
				Kv2MountPath: "secret",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := &VaultConfig{
				AuthMethod:   tt.fields.AuthMethod,
				VaultToken:   tt.fields.VaultToken,
				VaultAddr:    tt.fields.VaultAddr,
				SecretId:     tt.fields.SecretId,
				RoleId:       tt.fields.RoleId,
				PathPrefix:   tt.fields.PathPrefix,
				SecretIdFile: tt.fields.SecretIdFile,
				Kv2MountPath: tt.fields.Kv2MountPath,
			}
			if err := conf.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
