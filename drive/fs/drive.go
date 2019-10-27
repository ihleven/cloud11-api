package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/ihleven/cloud11-api/auth"
	"github.com/ihleven/cloud11-api/drive"
)

var Drive FSWebDrive = FSWebDrive{Root: "/Users/mi/tmp", Prefix: "/home", ServeURL: "/serve/home", PermissionMode: 0}

type FSWebDrive struct {
	// Absolute path inside filesystem
	Root string // /home/ihle/tmp
	// pathname of root dir in webview
	Prefix string // /home
	// pathname of root dir in serveview
	ServeURL string // /serve/home
	// indicates if index.html is served for directories in serveview
	serveIndexHtml bool
	// AlbumURL       string          `json:"albumUrl"`

	// Owner          *drive.User     `json:"-"` // alle Dateien gehören automatisch diesem User ( => homes )
	// Group          *drive.Group    `json:"-"` // jedes File des Storage bekommt automatisch diese Gruppe ( z.B. brunhilde )

	// PermMode overwrites file permissions globally (if set).
	// e.g., for a public readable but not writable filesystem: 0444
	PermissionMode os.FileMode              // wenn gesetzt erhält jedes File dies Permission =< wird nicht mehr auf fs gelesen
	Accounts       map[string]*auth.Account //
}

func (wd *FSWebDrive) GetServeHandle(path string) (os.FileInfo, *os.File, error) {

	location := strings.Replace(filepath.Clean(path), wd.ServeURL, wd.Root, 1)
	fmt.Println("GetServeHandle", path, location)
	info, err := os.Stat(location)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, errors.Wrapf(err, "os.Stat failed for %s (location: %s)", path, location)
		}
		return nil, nil, errors.Wrapf(err, "os.Stat failed for %s (location: %s)", path, location)
	}

	fd, err := os.OpenFile(location, os.O_RDONLY, 0)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "OpenFile failed for %s (location: %s)", path, location)
	}

	return info, fd, nil
}

func (wd *FSWebDrive) GetHandle(url, prefix string) (drive.Handle, error) {

	// path := strings.TrimPrefix(filepath.Clean(url), prefix)
	// location := filepath.Join(wd.Root, path)
	location := strings.Replace(url, prefix, wd.Root, 1)

	info, err := os.Stat(location)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.Wrapf(err, "os.Stat failed for %s (location: %s)", url, location)
		}
		return nil, errors.Wrapf(err, "os.Stat failed for %s (location: %s)", url, location)
	}

	fh := handle{FileInfo: info, location: location, mode: info.Mode()}

	if wd.PermissionMode != 0 {
		// replace 9 least significant bits from mode with storage.PermissionMode
		fh.mode = (fh.mode & 0xfffffe00) | (wd.PermissionMode & os.ModePerm) // & 0x1ff
	}
	return &fh, nil
}

func (wd *FSWebDrive) GetFile(url string, account *auth.Account) (*drive.File, error) {

	hndl, err := wd.GetHandle(url, wd.Prefix)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetHandle(%v, %v)", url, wd.Prefix)
	}
	//location := strings.Replace(filepath.Clean(url), wd.Prefix, wd.Root, 1)

	fmt.Printf(" * GetFile(%v) -> %v, %v\n", url, hndl.(*handle).location, err)

	fh := hndl.(*handle)
	uid, gid := fh.getUidGid()
	file := drive.File{
		Handle:        hndl,
		URL:           url,
		Name:          hndl.Name(),
		Size:          hndl.Size(),
		Mode:          fh.mode,
		Type:          fh.GuessMIME(),
		Permissions:   fh.mode.String(),
		Owner:         GetUserByID(uid),
		Group:         GetGroupByID(gid),
		Authorization: fh.GetPermissions(account),
		Modified:      hndl.ModTime(),
	}

	return &file, nil
}

func (wd *FSWebDrive) GetFolder(file *drive.File, account *auth.Account) (*drive.Folder, error) {

	handles, err := file.ReadDir(wd.PermissionMode)
	if err != nil {
		return nil, err
	}

	folder := drive.Folder{File: file, Entries: make([]drive.File, len(handles))}

	for index, handle := range handles {

		folder.Entries[index] = drive.File{
			Handle: handle,
			URL:    filepath.Join(file.URL, handle.Name()),
			Name:   handle.Name(),
			Size:   handle.Size(),
		}
	}
	return &folder, nil
}
