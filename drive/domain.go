package drive

import (
	"io"
	"os"
	"time"

	"github.com/ihleven/cloud11-api/auth"
)

// WebDrive is the domain and can modify state, interacting with storage and/or manipulating data as needed.
// It contains the business logic.

type Driver interface {
	Open(string) (Handle, error)
	OpenFile(string, *auth.Account) (*File, error)
	Create(string) (Handle, error)
	Mkdir(string) (Handle, error)
	ListFiles(*File, *auth.Account) ([]File, error)
	// CreateFile(folder *File, name string) (Handle, error)
	// Mkdir(string) (Handle, error)
}

type Handle interface {
	//Name() string       // base name of the file
	//Size() int64        // length in bytes for regular files; system-dependent for others
	//Mode() os.FileMode  // file mode bits
	//ModTime() time.Time // modification time
	IsDir() bool // abbreviation for Mode().IsDir()
	//Sys() interface{}   // underlying data source (can return nil)
	OpenFile(flag int, perm os.FileMode) (*os.File, error)
	ReadDir(mode os.FileMode) ([]Handle, error)
	io.Reader
	io.Writer
	//io.Seeker
	//io.Closer

	HasReadPermission(*auth.Account) bool
}

// //FileResponder builds the entire HTTP response from the domain's output which is given to it by the action.
// type FileResponder struct {
// 	handle Handle
// }

// File bundles all publically available information about Files (and Folders).
type File struct {
	Handle        `json:"-"`
	URL           string        `json:"url"`
	Name          string        `json:"name"`
	Size          int64         `json:"size"`
	Mode          os.FileMode   `json:"mode"`
	Type          Type          `json:"type"`
	Permissions   string        `json:"permissions"`
	Owner         *User         `json:"owner"`
	Group         *Group        `json:"group"`
	Authorization Authorization `json:"auth"`
	//Created     *time.Time   `json:"created"`
	Modified time.Time `json:"modified"`
	//Accessed    *time.Time   `json:"accessed"`
}

type Type struct {
	Filetype  string `json:"filetype"`
	Mediatype string `json:"mediatype"`
	Subtype   string `json:"subtype"`
	MIME      string `json:"mime"`
	Charset   string `json:"charset"`
}

type User struct {
	Uid      string `json:"uid"`
	Gid      string `json:"-"`
	Username string `json:"name"`
	Name     string `json:"-"`
	HomeDir  string `json:"-"`
}

type Group struct {
	Gid  string `json:"gid"`  // group ID
	Name string `json:"name"` // group name
}

type Authorization struct {
	IsOwner bool `json:"isOwner"`
	InGroup bool `json:"inGroup"`
	R       bool `json:"read"`
	W       bool `json:"write"`
	X       bool `json:"exec"`
}

type Folder struct {
	*File
	Account     *auth.Account `json:"account"`
	Drive       Driver        `json:"drive"`
	Breadcrumbs []Breadcrumb  `json:"breadcrumbs"`
	Entries     []File        `json:"entries"`
}

type Breadcrumb struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type DriveAction struct {
	//Template string
	Account     *auth.Account `json:"account"`
	Drive       Driver        `json:"drive"`
	Breadcrumbs []Breadcrumb  `json:"breadcrumbs"`
	*File
	Content string `json:"content,omitempty"`
	Entries []File `json:"entries,omitempty"`
	Image   string `json:"image,omitempty"`
	path    string
}
