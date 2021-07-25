package log

import (
	"fmt"
	"testing"
)

func Test_filesize(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "1b", args: "1b", want: 1, wantErr: false},
		{name: "1kb", args: "1kb", want: 1 * kb, wantErr: false},
		{name: "1mb", args: "1mb", want: 1 * mb, wantErr: false},
		{name: "1gb", args: "1gb", want: 1 * gb, wantErr: false},
		{name: "1", args: "1", want: 1, wantErr: false},
		{name: "1tb", args: "1tb", want: 0, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := filesize(tt.args)
			fmt.Println(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("filesize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("filesize() got = %v, want %v", got, tt.want)
			}
		})
	}
}
