package config

import "testing"

func TestFsWriterConfig_Validate(t *testing.T) {
	type fields struct {
		PrivateKeyPath string
		CertPath       string
		PemPath        string
		Username       string
		Group          string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid - only pem",
			fields: fields{
				PrivateKeyPath: "",
				CertPath:       "",
				PemPath:        "/path/to/pem",
				Username:       "root",
				Group:          "root",
			},
			wantErr: false,
		},
		{
			name: "valid - private key, cert and pem",
			fields: fields{
				PrivateKeyPath: "/path/to/private",
				CertPath:       "/path/to/cert",
				PemPath:        "/path/to/pem",
				Username:       "root",
				Group:          "root",
			},
			wantErr: false,
		},
		{
			name: "valid - no pem",
			fields: fields{
				PrivateKeyPath: "/path/to/private",
				CertPath:       "/path/to/cert",
				Username:       "root",
				Group:          "root",
			},
			wantErr: false,
		},
		{
			name: "invalid - no private key but pem",
			fields: fields{
				CertPath: "/path/to/cert",
				PemPath:  "/path/to/pem",
				Username: "root",
				Group:    "root",
			},
			wantErr: true,
		},
		{
			name: "invalid - nada",
			fields: fields{
				Username: "root",
				Group:    "root",
			},
			wantErr: true,
		},
		{
			name: "invalid - no user name",
			fields: fields{
				PrivateKeyPath: "/path/to/private",
				CertPath:       "/path/to/cert",
				Group:          "root",
			},
			wantErr: true,
		},
		{
			name: "invalid - no group name",
			fields: fields{
				PrivateKeyPath: "/path/to/private",
				CertPath:       "/path/to/cert",
				Username:       "root",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := FsWriterConfig{
				PrivateKeyPath: tt.fields.PrivateKeyPath,
				CertPath:       tt.fields.CertPath,
				PemPath:        tt.fields.PemPath,
				Username:       tt.fields.Username,
				Group:          tt.fields.Group,
			}
			if err := conf.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
