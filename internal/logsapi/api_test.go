package logsapi

import (
	"reflect"
	"strconv"
	"testing"
)

func TestHandleFilePagination(t *testing.T) {
	logData := []byte("abcdefghijklmnopqrstuvwxyz")
	fileSize := len(logData)

	tests := []struct {
		name        string
		offsetParam string
		limitParam  string
		want        []byte
	}{
		{
			name:        "Valid offset and limit",
			offsetParam: "2",
			limitParam:  "5",
			want:        logData[2:7],
		},
		{
			name:        "Offset zero, limit zero (should return all)",
			offsetParam: "0",
			limitParam:  "0",
			want:        logData,
		},
		{
			name:        "Offset negative, limit positive (should treat offset as 0)",
			offsetParam: "-5",
			limitParam:  "4",
			want:        logData[:4],
		},
		{
			name:        "Start > end (should return empty)",
			offsetParam: strconv.Itoa(fileSize - 1),
			limitParam:  "-100",
			want:        logData[fileSize-1:],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handleFilePagination(tt.offsetParam, tt.limitParam, fileSize, logData)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleFilePagination(%q, %q, %d, logData) = %q, want %q",
					tt.offsetParam, tt.limitParam, fileSize, got, tt.want)
			}
		})
	}
}
