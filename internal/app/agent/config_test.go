package agent

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Imomali1/metrics/internal/pkg/utils"
)

func TestParse(t *testing.T) {
	cfg := LoadConfig()
	require.NotEmpty(t, cfg)
}

func Test_getEnvInt(t *testing.T) {
	type args struct {
		key           string
		flagValue     int
		fileConfValue *int
		defaultValue  int
	}

	tests := []struct {
		name     string
		initFunc func()
		args     args
		want     int
	}{
		{
			name:     "with valid os env",
			initFunc: func() { os.Setenv("NUMBER", "250") },
			args: args{
				key:           "NUMBER",
				flagValue:     125,
				fileConfValue: utils.Ptr(150),
				defaultValue:  0,
			},
			want: 250,
		},
		{
			name:     "with invalid os env",
			initFunc: func() { os.Setenv("NUMBER", "invalid") },
			args: args{
				key:           "NUMBER",
				flagValue:     125,
				fileConfValue: utils.Ptr(150),
				defaultValue:  0,
			},
			want: 125,
		},
		{
			name:     "from flag",
			initFunc: func() { os.Unsetenv("NUMBER") },
			args: args{
				key:           "NUMBER",
				flagValue:     125,
				fileConfValue: utils.Ptr(150),
				defaultValue:  0,
			},
			want: 125,
		},
		{
			name:     "from file",
			initFunc: func() { os.Unsetenv("NUMBER") },
			args: args{
				key:           "NUMBER",
				flagValue:     0,
				fileConfValue: utils.Ptr(150),
				defaultValue:  0,
			},
			want: 150,
		},
		{
			name:     "default",
			initFunc: func() { os.Unsetenv("NUMBER") },
			args: args{
				key:           "NUMBER",
				flagValue:     125,
				fileConfValue: nil,
				defaultValue:  125,
			},
			want: 125,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.initFunc()
			if got := getEnvInt(tt.args.key, tt.args.flagValue, tt.args.fileConfValue, tt.args.defaultValue); got != tt.want {
				t.Errorf("getEnvInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getEnvString(t *testing.T) {
	type args struct {
		key           string
		flagValue     string
		fileConfValue *string
		defaultValue  string
	}

	tests := []struct {
		name     string
		initFunc func()
		args     args
		want     string
	}{
		{
			name:     "normal os env",
			initFunc: func() { os.Setenv("TEXT", "normal") },
			args: args{
				key:       "TEXT",
				flagValue: "argument",
			},
			want: "normal",
		},
		{
			name:     "empty os env",
			initFunc: func() { os.Setenv("TEXT", "") },
			args: args{
				key:       "TEXT",
				flagValue: "argument",
			},
			want: "",
		},
		{
			name:     "from flag",
			initFunc: func() { os.Unsetenv("TEXT") },
			args: args{
				key:       "TEXT",
				flagValue: "flag",
			},
			want: "flag",
		},
		{
			name:     "from file",
			initFunc: func() { os.Unsetenv("TEXT") },
			args: args{
				key:           "TEXT",
				fileConfValue: utils.Ptr("file"),
			},
			want: "file",
		},
		{
			name:     "default",
			initFunc: func() { os.Unsetenv("TEXT") },
			args: args{
				key:          "TEXT",
				flagValue:    "default",
				defaultValue: "default",
			},
			want: "default",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.initFunc()
			if got := getEnvString(tt.args.key, tt.args.flagValue, tt.args.fileConfValue, tt.args.defaultValue); got != tt.want {
				t.Errorf("getEnvString() = %v, want %v", got, tt.want)
			}
		})
	}
}
