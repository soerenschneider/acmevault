package config

import (
	"reflect"
	"testing"
)

func TestFromFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    AcmeVaultClientConfig
		wantErr bool
	}{
		{
			name: "Valid file",
			args: args{
				path: "../../contrib/client.json",
			},
			want: AcmeVaultClientConfig{
				FsWriterConfig: FsWriterConfig{
					PrivateKeyFile: "/tmp/mydomain.key",
					CertFile:       "/tmp/mydomain.crt",
					Username:       "myusername",
					Group:          "mygroup",
				},
				VaultConfig: VaultConfig{
					RoleId:     "roleId",
					SecretId:   "secretId",
					VaultAddr:  "https://vault:8200",
					PathPrefix: "my-prefix",
				},
			},
			wantErr: false,
		},
		{
			name: "Non existing file",
			args: args{
				path: "../../contrib/nonexistent.json",
			},
			want:    AcmeVaultClientConfig{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AcmeVaultClientConfigFromFile(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}
