package tar_utils

import (
	"io"
	"testing"
)

func TestPath2TarReader(t *testing.T) {
	type args struct {
		sourcePath string
	}
	tests := []struct {
		name    string
		args    args
		want    io.Reader
		wantErr bool
	}{
		{
			name: "Normal",
			args: args{
				sourcePath: "../tmp/tar_utils_common_test_normal",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Path2TarReader(tt.args.sourcePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Path2TarReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
