package config

import "testing"

func TestVaultConfig_Validate(t *testing.T) {
	type fields struct {
		VaultToken            string
		VaultAddr             string
		SecretId              string
		RoleId                string
		TokenIncreaseSeconds  int
		TokenIncreaseInterval int
		PathPrefix            string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid config - token",
			fields: fields{
				VaultToken:            "s.asd83hrfhasfjsda",
				VaultAddr:             "https://my-vault-instance:443",
				TokenIncreaseSeconds:  0,
				TokenIncreaseInterval: 0,
				PathPrefix:            "production",
			},
		},
		{
			name: "valid config - approle",
			fields: fields{
				VaultAddr:             "https://my-vault-instance:443",
				SecretId:              "super-secret",
				RoleId:                "my-role",
				TokenIncreaseSeconds:  0,
				TokenIncreaseInterval: 0,
				PathPrefix:            "dev-v002",
			},
		},
		{
			name: "invalid config - missing protocol",
			fields: fields{
				VaultToken:            "s.asd83hrfhasfjsda",
				VaultAddr:             "my-vault-instance:443",
				TokenIncreaseSeconds:  0,
				TokenIncreaseInterval: 0,
				PathPrefix:            "production",
			},
			wantErr: true,
		},
		{
			name: "invalid config - missing port",
			fields: fields{
				VaultToken:            "s.asd83hrfhasfjsda",
				VaultAddr:             "http://my-vault-instance",
				TokenIncreaseSeconds:  0,
				TokenIncreaseInterval: 0,
				PathPrefix:            "production",
			},
			wantErr: true,
		},
		{
			name: "invalid config - invalid path prefix",
			fields: fields{
				VaultToken:            "s.asd83hrfhasfjsda",
				VaultAddr:             "http://my-vault-instance:443",
				TokenIncreaseSeconds:  0,
				TokenIncreaseInterval: 0,
				PathPrefix:            "/production",
			},
			wantErr: true,
		},
		{
			name: "invalid config - no auth methods",
			fields: fields{
				VaultAddr:             "http://my-vault-instance:443",
				TokenIncreaseSeconds:  0,
				TokenIncreaseInterval: 0,
				PathPrefix:            "production",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := &VaultConfig{
				VaultToken:            tt.fields.VaultToken,
				VaultAddr:             tt.fields.VaultAddr,
				SecretId:              tt.fields.SecretId,
				RoleId:                tt.fields.RoleId,
				TokenIncreaseSeconds:  tt.fields.TokenIncreaseSeconds,
				TokenIncreaseInterval: tt.fields.TokenIncreaseInterval,
				PathPrefix:            tt.fields.PathPrefix,
			}
			if err := conf.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
