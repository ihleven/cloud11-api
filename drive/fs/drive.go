package fs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
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

func (wd *FSWebDrive) GetFiles(file *drive.File, account *auth.Account) ([]drive.File, error) {

	handles, err := file.ReadDir(wd.PermissionMode)
	if err != nil {
		return nil, err
	}

	entries := make([]drive.File, len(handles))

	for index, Handle := range handles {
		fh := Handle.(*handle)
		uid, gid := fh.getUidGid()
		entries[index] = drive.File{
			Handle:        fh,
			URL:           filepath.Join(file.URL, Handle.Name()),
			Name:          Handle.Name(),
			Size:          Handle.Size(),
			Mode:          fh.mode,
			Type:          fh.GuessMIME(),
			Permissions:   fh.mode.String(),
			Owner:         GetUserByID(uid),
			Group:         GetGroupByID(gid),
			Authorization: fh.GetPermissions(account),
			Modified:      Handle.ModTime(),
		}
	}
	return entries, nil
}

func (wd *FSWebDrive) GetFolder(file *drive.File, account *auth.Account) (*drive.Folder, error) {

	handles, err := file.ReadDir(wd.PermissionMode)
	if err != nil {
		return nil, err
	}

	folder := drive.Folder{File: file, Entries: make([]drive.File, len(handles))}

	for index, Handle := range handles {
		fh := Handle.(*handle)
		uid, gid := fh.getUidGid()
		folder.Entries[index] = drive.File{
			Handle:        fh,
			URL:           filepath.Join(file.URL, Handle.Name()),
			Name:          Handle.Name(),
			Size:          Handle.Size(),
			Mode:          fh.mode,
			Type:          fh.GuessMIME(),
			Permissions:   fh.mode.String(),
			Owner:         GetUserByID(uid),
			Group:         GetGroupByID(gid),
			Authorization: fh.GetPermissions(account),
			Modified:      Handle.ModTime(),
		}
	}
	return &folder, nil
}
func (wd *FSWebDrive) CreateFile(folder *drive.File, name string) (drive.Handle, error) {
	h := folder.Handle.(*handle)
	l := path.Join(h.location, name)
	//
	var _, err = os.Stat(l)
	var file *os.File
	// create file if not exists
	if os.IsNotExist(err) {
		file, err = os.Create(l)

		//defer file.Close()
	} else {
		basename := strings.TrimSuffix(name, filepath.Ext(name)) + ".*" + filepath.Ext(name)
		file, err = ioutil.TempFile(h.location, basename)
	}
	// file, err := os.OpenFile(h.location+"/"+name, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not create file %v", name)
	}
	info, err := file.Stat()
	handle := NewHandle(info, path.Join(h.location, info.Name()), 0)
	return handle, nil
}

func (wd *FSWebDrive) Mkdir(url string) (drive.Handle, error) {

	location := strings.Replace(path.Clean(url), wd.Prefix, wd.Root, 1)

	err := os.Mkdir(location, 0)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not create folder %v", url)
	}
	handle, err := wd.GetHandle(path.Join(wd.Root, url), "")

	return handle, nil
}

func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}
func (wd *FSWebDrive) GenerateBreadcrumbs(p string) (breadcrumbs []drive.Breadcrumb) {
	var url = "/"
	for elem, remainder := ShiftPath(p); elem != ""; elem, remainder = ShiftPath(remainder) {
		url = path.Join(url, elem)
		bc := drive.Breadcrumb{Name: elem, URL: url}
		breadcrumbs = append(breadcrumbs, bc)
	}
	return
}

func (wd *FSWebDrive) GenerateParents(p string) []drive.Breadcrumb {
	var path string
	elements := strings.Split(p[1:], "/")
	list := make([]drive.Breadcrumb, len(elements)+1)
	list[0] = drive.Breadcrumb{Name: "cloud11", URL: "/"}
	for index, element := range elements {
		path = fmt.Sprintf("%s/%s", path, element)
		list[index+1] = drive.Breadcrumb{Name: element, URL: path}
	}
	return list
}
