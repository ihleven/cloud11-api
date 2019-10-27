package fs

import (
	"testing"
)

func equal(a, b *FSHandle) bool {
	return a.drive == b.drive && a.name == b.name && a.isDir == b.isDir
}
func TestGetHandle(t *testing.T) {

	var drive = FSWebDrive{Root: "/Users/mi/tmp/14", Prefix: "home"}

	var tests = []struct {
		input string
		want  FSHandle
	}{
		{"home/DSC02316.jpg", FSHandle{drive: &drive, name: "DSC02316.jpg", isDir: false}},
	}
	for _, test := range tests {
		t.Log(drive.GetHandle(test.input))
		if got, err := drive.GetHandle(test.input); !equal(got, &test.want) {
			t.Log(got, err)
			t.Errorf("parseURL(%q) = %v, want: %v", test.input, *got, test.want)
		}
	}
}
