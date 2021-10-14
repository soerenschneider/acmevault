package config

import "testing"

func TestFsWriterConfig_Validate(t *testing.T) {
	type fields struct {
		PrivateKeyFile string
		CertFile       string
		PemFile        string
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
				PrivateKeyFile: "",
				CertFile:       "",
				PemFile:        "/path/to/pem",
				Username:       "root",
				Group:          "root",
			},
			wantErr: false,
		},
		{
			name: "valid - private key, cert and pem",
			fields: fields{
				PrivateKeyFile: "/path/to/private",
				CertFile:       "/path/to/cert",
				PemFile:        "/path/to/pem",
				Username:       "root",
				Group:          "root",
			},
			wantErr: false,
		},
		{
			name: "valid - no pem",
			fields: fields{
				PrivateKeyFile: "/path/to/private",
				CertFile:       "/path/to/cert",
				Username:       "root",
				Group:          "root",
			},
			wantErr: false,
		},
		{
			name: "valid - no private key but pem",
			fields: fields{
				CertFile: "/path/to/cert",
				PemFile:  "/path/to/pem",
				Username: "root",
				Group:    "root",
			},
			wantErr: false,
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
				PrivateKeyFile: "/path/to/private",
				CertFile:       "/path/to/cert",
				Group:          "root",
			},
			wantErr: true,
		},
		{
			name: "invalid - no group name",
			fields: fields{
				PrivateKeyFile: "/path/to/private",
				CertFile:       "/path/to/cert",
				Username:       "root",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := FsWriterConfig{
				PrivateKeyFile: tt.fields.PrivateKeyFile,
				CertFile:       tt.fields.CertFile,
				PemFile:        tt.fields.PemFile,
				Username:       tt.fields.Username,
				Group:          tt.fields.Group,
			}
			if err := conf.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
