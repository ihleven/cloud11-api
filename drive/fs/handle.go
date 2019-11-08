package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"syscall"
	"time"

	"github.com/ihleven/cloud11-api/auth"
	"github.com/ihleven/cloud11-api/drive"
	"github.com/pkg/errors"
)

// FSHandle erfüllt Mimer, Locator,
type FSHandle struct {
	drive   *FSWebDrive
	Path    string      // Pfad relativ zur Storage-Wurzel
	Name    string      // base name of the file
	Size    int64       // length in bytes for regular files; system-dependent for others
	Mode    os.FileMode // file mode bits
	ModTime time.Time   // modification time
	IsDir   bool        // abbreviation for Mode().IsDir()
	//Sys() interface{}   // underlying data source (can return nil)
}

func (h *FSHandle) Open() (*os.File, error) {

	location := filepath.Join(h.drive.Root, h.Path)
	return os.OpenFile(location, os.O_RDONLY, 0)
	// os.OpenFile(name string, flag int, perm FileMode) (*File, error) {
}

func NewHandle(fileInfo os.FileInfo, location string, mode os.FileMode) *handle {

	handle := &handle{
		FileInfo: fileInfo,
		location: location,
		mode:     fileInfo.Mode(),
	}
	if mode != 0 {
		// replace 9 least significant bits from mode with storage.PermissionMode
		handle.mode = (handle.mode & 0xfffffe00) | (mode & os.ModePerm) // & 0x1ff
	}
	return handle
}

type handle struct {
	os.FileInfo
	location string
	mode     os.FileMode
}

func (h handle) Mode() os.FileMode {
	return h.mode
}

func (h handle) OpenFile(flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(h.location, flag, perm)
}

func (fh handle) Read(b []byte) (n int, err error) {

	fd, err := fh.OpenFile(0, 0)
	if err != nil {
		return 0, err
	}
	defer fd.Close()
	fd.Seek(0, 0)

	bytes, err := fd.Read(b)
	if err != nil {
		return bytes, err
	}
	if bytes != int(fh.Size()) {
		return bytes, errors.Errorf("read only %d of %d bytes", bytes, fh.Size())
	}
	return bytes, nil
}

func (fh handle) Write(b []byte) (int, error) {
	fd, err := fh.OpenFile(os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0)
	if err != nil {
		return 0, err
	}
	defer fd.Close()

	n, err := fd.Write(b)
	if err != nil {
		return n, errors.Wrapf(err, "Could not write to handle %v: %v", fh.Name(), string(b))
	}
	return n, nil
}

func (h handle) ReadDir(mode os.FileMode) ([]drive.Handle, error) {

	fd, err := os.Open(h.location)
	defer fd.Close()
	if err != nil {
		return nil, errors.Wrapf(err, "ReadDir: Could not get file descriptor for %v", h.location)
	}

	files, err := fd.Readdir(-1)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not read dir %v", fd)
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })

	entries := make([]drive.Handle, len(files))
	for index, info := range files {

		entries[index] = NewHandle(info, filepath.Join(h.location, info.Name()), mode)
	}
	return entries, nil
}

func (fh handle) getUidGid() (uid uint32, gid uint32) {

	if stat, ok := fh.Sys().(*syscall.Stat_t); ok {
		uid, gid = stat.Uid, stat.Gid
	}
	return
}

//PERMISSIONS

func (fh handle) HasReadPermission(account *auth.Account) bool {

	if fh.mode&OS_OTH_R != 0 {
		fmt.Println("has other read:", fh.mode, fh.mode&4)
		return true
	}

	stat, ok := fh.Sys().(*syscall.Stat_t)
	if ok {
		fmt.Println("checking group:", account, fh.mode&OS_GROUP_R)
		if fh.mode&OS_GROUP_R != 0 {
			return account != nil && stat.Gid == gid[account.Username]
		}
		if fh.mode&OS_USER_R != 0 {
			return account != nil && stat.Uid == uid[account.Username]
		}
	}
	return false
}

func (fh handle) GetPermissions(account *auth.Account) drive.Authorization { // => handle

	perm := drive.Authorization{}

	if account != nil {
		if stat, ok := fh.Sys().(*syscall.Stat_t); ok {
			perm.InGroup = stat.Gid == gid[account.Username]
			perm.IsOwner = stat.Uid == uid[account.Username]
		}
	}

	read, write, x := OS_OTH_R, OS_OTH_W, OS_OTH_X
	if perm.InGroup {
		read, write, x = read|OS_GROUP_R, write|OS_GROUP_W, x|OS_GROUP_X
	}
	if perm.IsOwner {
		read, write, x = read|OS_USER_R, write|OS_USER_W, x|OS_USER_X
	}

	perm.R = fh.mode&os.FileMode(read) != 0
	perm.W = fh.mode&os.FileMode(write) != 0
	perm.X = fh.mode&os.FileMode(x) != 0
	return perm
}

const (
	OS_READ        = 04
	OS_WRITE       = 02
	OS_EX          = 01
	OS_USER_SHIFT  = 6
	OS_GROUP_SHIFT = 3
	OS_OTH_SHIFT   = 0

	OS_USER_R   = OS_READ << OS_USER_SHIFT
	OS_USER_W   = OS_WRITE << OS_USER_SHIFT
	OS_USER_X   = OS_EX << OS_USER_SHIFT
	OS_USER_RW  = OS_USER_R | OS_USER_W
	OS_USER_RWX = OS_USER_RW | OS_USER_X

	OS_GROUP_R   = OS_READ << OS_GROUP_SHIFT
	OS_GROUP_W   = OS_WRITE << OS_GROUP_SHIFT
	OS_GROUP_X   = OS_EX << OS_GROUP_SHIFT
	OS_GROUP_RW  = OS_GROUP_R | OS_GROUP_W
	OS_GROUP_RWX = OS_GROUP_RW | OS_GROUP_X

	OS_OTH_R   = OS_READ << OS_OTH_SHIFT
	OS_OTH_W   = OS_WRITE << OS_OTH_SHIFT
	OS_OTH_X   = OS_EX << OS_OTH_SHIFT
	OS_OTH_RW  = OS_OTH_R | OS_OTH_W
	OS_OTH_RWX = OS_OTH_RW | OS_OTH_X

	OS_ALL_R   = OS_USER_R | OS_GROUP_R | OS_OTH_R
	OS_ALL_W   = OS_USER_W | OS_GROUP_W | OS_OTH_W
	OS_ALL_X   = OS_USER_X | OS_GROUP_X | OS_OTH_X
	OS_ALL_RW  = OS_ALL_R | OS_ALL_W
	OS_ALL_RWX = OS_ALL_RW | OS_GROUP_X
)
