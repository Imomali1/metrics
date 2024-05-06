package storage

import (
	"context"
	"reflect"
	"testing"
)

func TestNewStorage(t *testing.T) {
	type args struct {
		opts []OptionsStorage
	}
	tests := []struct {
		name    string
		args    args
		want    *Storage
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewStorage(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewStorage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStorage() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRestoreFile(t *testing.T) {
	type args struct {
		ctx      context.Context
		filename string
	}
	tests := []struct {
		name string
		args args
		want OptionsStorage
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RestoreFile(tt.args.ctx, tt.args.filename); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RestoreFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDB(t *testing.T) {
	type args struct {
		ctx context.Context
		dsn string
	}
	tests := []struct {
		name string
		args args
		want OptionsStorage
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDB(tt.args.ctx, tt.args.dsn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithSyncWrite(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want OptionsStorage
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithSyncWrite(tt.args.filename); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithSyncWrite() = %v, want %v", got, tt.want)
			}
		})
	}
}
