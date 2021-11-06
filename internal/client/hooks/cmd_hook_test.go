package hooks

import "testing"

func TestOsCommandPostHook_Invoke(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "cmd",
			args:    []string{"date"},
			wantErr: false,
		},
		{
			name:    "cmd with arg",
			args:    []string{"date", "+%s"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hook := &CommandPostHook{
				commands: tt.args,
			}
			if err := hook.Invoke(); (err != nil) != tt.wantErr {
				t.Errorf("Invoke() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
