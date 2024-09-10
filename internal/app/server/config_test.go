package server

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Imomali1/metrics/internal/pkg/utils"
)

func TestLoadConfig(t *testing.T) {
	cfg := LoadConfig()
	require.NotEmpty(t, cfg)
}

func TestLoadFileConfig(t *testing.T) {
	// invalid path
	cfg, err := LoadFileConfig("/invalid/path/to/config/file")
	require.Error(t, err)
	require.Empty(t, cfg)

	// empty file
	file, err := os.CreateTemp("", "config-*.json")
	require.NoError(t, err)
	defer os.Remove(file.Name())
	defer file.Close()
	cfg, err = LoadFileConfig(file.Name())
	require.Error(t, err)
	require.Empty(t, cfg)

	// write normal json, but not FileConfig struct
	_, err = file.Write([]byte(`{"invalid-field":"invalid-value"}`))
	require.NoError(t, err)
	cfg, err = LoadFileConfig(file.Name())
	require.NoError(t, err)
	require.Empty(t, cfg)

	// write normal json with struct FileConfig
	file2, err := os.CreateTemp("", "config-*.json")
	require.NoError(t, err)
	defer os.Remove(file2.Name())
	defer file2.Close()
	_, err = file2.Write([]byte(`{"address":"localhost:8080"}`))
	require.NoError(t, err)
	cfg, err = LoadFileConfig(file2.Name())
	require.NoError(t, err)
	require.NotEmpty(t, cfg)
	require.Equal(t, *cfg.ServerAddress, "localhost:8080")
}

func Test_getEnvBool(t *testing.T) {
	type args struct {
		key          string
		flagValue    bool
		fileValue    *bool
		defaultValue bool
	}

	tests := []struct {
		name     string
		initFunc func()
		args     args
		want     bool
	}{
		{
			name:     "with normal os env",
			initFunc: func() { os.Setenv("ENV_KEY", "true") },
			args: args{
				key:          "ENV_KEY",
				flagValue:    false,
				fileValue:    utils.Ptr(false),
				defaultValue: false,
			},
			want: true,
		},
		{
			name:     "with empty os env",
			initFunc: func() { os.Setenv("ENV_KEY", "") },
			args: args{
				key:          "ENV_KEY",
				flagValue:    true,
				fileValue:    utils.Ptr(false),
				defaultValue: false,
			},
			want: true,
		},
		{
			name:     "from flag",
			initFunc: func() { os.Unsetenv("ENV_KEY") },
			args: args{
				key:          "ENV_KEY",
				flagValue:    true,
				fileValue:    utils.Ptr(false),
				defaultValue: false,
			},
			want: true,
		},
		{
			name:     "from file",
			initFunc: func() { os.Unsetenv("ENV_KEY") },
			args: args{
				key:          "ENV_KEY",
				flagValue:    false,
				fileValue:    utils.Ptr(true),
				defaultValue: false,
			},
			want: true,
		},
		{
			name:     "default value",
			initFunc: func() { os.Unsetenv("ENV_KEY") },
			args: args{
				key:          "ENV_KEY",
				flagValue:    false,
				fileValue:    nil,
				defaultValue: true,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.initFunc()
			if got := getEnvBool(
				tt.args.key,
				tt.args.flagValue,
				tt.args.fileValue,
				tt.args.defaultValue,
			); got != tt.want {
				t.Errorf("getEnvBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getEnvInt(t *testing.T) {
	type args struct {
		key          string
		flagValue    int
		fileValue    *int
		defaultValue int
	}

	tests := []struct {
		name     string
		initFunc func()
		args     args
		want     int
	}{
		{
			name:     "with normal os env",
			initFunc: func() { os.Setenv("ENV_KEY", "1") },
			args: args{
				key:          "ENV_KEY",
				flagValue:    2,
				fileValue:    utils.Ptr(3),
				defaultValue: 4,
			},
			want: 1,
		},
		{
			name:     "with empty os env",
			initFunc: func() { os.Setenv("ENV_KEY", "") },
			args: args{
				key:          "ENV_KEY",
				flagValue:    2,
				fileValue:    utils.Ptr(3),
				defaultValue: 4,
			},
			want: 2,
		},
		{
			name:     "from flag",
			initFunc: func() { os.Unsetenv("ENV_KEY") },
			args: args{
				key:          "ENV_KEY",
				flagValue:    2,
				fileValue:    utils.Ptr(3),
				defaultValue: 4,
			},
			want: 2,
		},
		{
			name:     "from file",
			initFunc: func() { os.Unsetenv("ENV_KEY") },
			args: args{
				key:          "ENV_KEY",
				flagValue:    0,
				fileValue:    utils.Ptr(3),
				defaultValue: 4,
			},
			want: 3,
		},
		{
			name:     "default value",
			initFunc: func() { os.Unsetenv("ENV_KEY") },
			args: args{
				key:          "ENV_KEY",
				flagValue:    0,
				fileValue:    nil,
				defaultValue: 4,
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.initFunc()
			if got := getEnvInt(
				tt.args.key,
				tt.args.flagValue,
				tt.args.fileValue,
				tt.args.defaultValue,
			); got != tt.want {
				t.Errorf("getEnvInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getEnvString(t *testing.T) {
	type args struct {
		key          string
		flagValue    string
		fileValue    *string
		defaultValue string
	}

	tests := []struct {
		name     string
		initFunc func()
		args     args
		want     string
	}{
		{
			name:     "with normal os env",
			initFunc: func() { os.Setenv("ENV_KEY", "env") },
			args: args{
				key:          "ENV_KEY",
				flagValue:    "flag",
				fileValue:    utils.Ptr("file"),
				defaultValue: "default",
			},
			want: "env",
		},
		{
			name:     "with empty os env",
			initFunc: func() { os.Setenv("ENV_KEY", "") },
			args: args{
				key:          "ENV_KEY",
				flagValue:    "flag",
				fileValue:    utils.Ptr("file"),
				defaultValue: "default",
			},
			want: "",
		},
		{
			name:     "from flag",
			initFunc: func() { os.Unsetenv("ENV_KEY") },
			args: args{
				key:          "ENV_KEY",
				flagValue:    "flag",
				fileValue:    utils.Ptr("file"),
				defaultValue: "default",
			},
			want: "flag",
		},
		{
			name:     "from file",
			initFunc: func() { os.Unsetenv("ENV_KEY") },
			args: args{
				key:          "ENV_KEY",
				flagValue:    "",
				fileValue:    utils.Ptr("file"),
				defaultValue: "default",
			},
			want: "file",
		},
		{
			name:     "default value",
			initFunc: func() { os.Unsetenv("ENV_KEY") },
			args: args{
				key:          "ENV_KEY",
				flagValue:    "",
				fileValue:    nil,
				defaultValue: "default",
			},
			want: "default",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.initFunc()
			if got := getEnvString(
				tt.args.key,
				tt.args.flagValue,
				tt.args.fileValue,
				tt.args.defaultValue,
			); got != tt.want {
				t.Errorf("getEnvString() = %v, want %v", got, tt.want)
			}
		})
	}
}
