package hidrive

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/ihleven/cloud11-api/auth"
	"github.com/ihleven/cloud11-api/drive"
)

type hiHandle struct {
	// ctime,has_dirs,mtime,readable,size,type,writable
	Name     string `json:"name"`
	Path     string `json:"path"`
	Type     string `json:"type"`
	MIMEType string `json:"mime_type"`
	Size     uint64 `json:"size"`
	Readable bool   `json:"readable"`
	Writable bool   `json:"writable"`
	CTime    int64  `json:"ctime"`
	MTime    int64  `json:"mtime"`
	HasDirs  bool   `json:"has_dirs"`
	//Members  []Member `json:"members"`
	ID    string `json:"id"`
	Image Image  `json:"image"`
}

func (h *hiHandle) Mode() os.FileMode {
	return 0644
}
func (h *hiHandle) ModTime() time.Time {
	ut := time.Unix(h.MTime, 0)
	fmt.Println("ut:", ut)
	return ut
}
func (h *hiHandle) IsDir() bool {
	return h.Type == "dir"
}

func (h *hiHandle) OpenFile(flag int, perm os.FileMode) (*os.File, error) {
	return nil, nil
}
func (h *hiHandle) ReadDir(mode os.FileMode) ([]drive.Handle, error) {
	return nil, nil
}

func (h *hiHandle) HasReadPermission(*auth.Account) bool {
	return true
}
func (h hiHandle) Read(b []byte) (n int, err error) {
	return 0, nil
}
func (h hiHandle) Write(b []byte) (int, error) {
	return 0, nil
}

func (h hiHandle) ReadSeeker() (io.ReadSeeker, error) {

	buffer := make([]byte, h.Size)
	// read file content to buffer
	//file.Read(buffer)
	fileBytes := bytes.NewReader(buffer) // converted to io.ReadSeeker type
	return fileBytes, nil
}

func (h hiHandle) GuessMIME() drive.Type {

	var t = drive.Type{
		Filetype:  "",
		Mediatype: h.Type,
		Subtype:   "",
		MIME:      h.MIMEType,
		Charset:   "",
	}
	if h.Type == "dir" {
		t.Filetype = "D"
	}
	if h.Type == "file" {
		t.Filetype = "F"

		media, sub := path.Split(h.MIMEType)
		t.Mediatype = strings.TrimSuffix(media, "/")
		t.Subtype = sub
	}
	return t
}

func (h hiHandle) GetPermissions(account *auth.Account) drive.Authorization { // => handle

	perm := drive.Authorization{}

	perm.R = h.Readable
	perm.W = h.Writable
	perm.X = false
	return perm
}
