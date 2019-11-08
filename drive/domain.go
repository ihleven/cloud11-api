package drive

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/ihleven/cloud11-api/auth"
)

// WebDrive is the domain and can modify state, interacting with storage and/or manipulating data as needed.
// It contains the business logic.
type WebDrive interface { // WebDrive
	GetHandle(url string, prefix string) (Handle, error)
	GetFile(url string, account *auth.Account) (*File, error)
	GetFiles(*File, *auth.Account) ([]File, error)
	GetFolder(file *File, account *auth.Account) (*Folder, error)
	GenerateBreadcrumbs(p string) []Breadcrumb
	GenerateParents(p string) []Breadcrumb

	CreateFile(folder *File, name string) (Handle, error)
	Mkdir(string) (Handle, error)
}

type Handle interface {
	Name() string       // base name of the file
	Size() int64        // length in bytes for regular files; system-dependent for others
	Mode() os.FileMode  // file mode bits
	ModTime() time.Time // modification time
	//IsDir() bool        // abbreviation for Mode().IsDir()
	//Sys() interface{}   // underlying data source (can return nil)
	OpenFile(flag int, perm os.FileMode) (*os.File, error)
	ReadDir(mode os.FileMode) ([]Handle, error)
	io.Reader
	io.Writer
}

//FileResponder builds the entire HTTP response from the domain's output which is given to it by the action.
type FileResponder struct {
	handle Handle
}

type Filer interface {
}

// File bundles all publically available information about Files (and Folders).
//
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
	// Account *domain.Account `json:"-"`
	// Notation          string
	IsOwner bool
	InGroup bool
	R       bool
	W       bool
	X       bool
}

type Folder struct {
	*File
	Entries []File `json:"entries"`
}

type Breadcrumb struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Action takes HTTP requests (URLs and their methods)
// and uses that input to interact with the domain,
// after which it passes the domain's output to one and only one responder.
func Action(d WebDrive) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		//sessionUser, _ := session.GetSessionUser(r, w)
		var err error

		//if responder := GetActioneer(file, sessionUser); responder != nil {
		switch r.Method {
		case http.MethodGet:

		case http.MethodDelete:
		case http.MethodPost:

		case http.MethodPut:
		}
		if err != nil {

		}
		//}
	}
}
