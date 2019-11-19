package hidrive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/ihleven/cloud11-api/auth"
	"github.com/ihleven/cloud11-api/drive"
	"github.com/pkg/errors"
)

type hiDrive struct {
	// Absolute path inside filesystem
	Root     string `json:"-"`        // /users/ihleven
	URL      string `json:"url"`      // /hidrive
	ServeURL string `json:"serveUrl"` // /hiserve
	Token    Token  `json:"-"`
}

var HIDrive hiDrive = hiDrive{Root: "/users/ihleven", URL: "/hidrive", ServeURL: "/hiserve"}

func (hd *hiDrive) GetMeta(path string) (*hiHandle, error) {

	queryParams := url.Values{
		"path":   {path},
		"fields": {"name,path,ctime,has_dirs,mtime,readable,size,type,writable,rshare,zone,image.exif,image.width,image.height,mime_type,id"},
	}

	request, _ := http.NewRequest("GET", "https://api.hidrive.strato.com/2.1/meta", nil)
	request.URL.RawQuery = queryParams.Encode()
	request.Header.Set("Authorization", "Bearer "+hd.Token.GetAccessToken())

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, NewHiDriveError(res.Body, res.StatusCode, res.Status)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println("hidriveGetMetaResponse:", string(body), res.Status, res.StatusCode, "adf")

	var handle hiHandle
	if err = json.NewDecoder(bytes.NewReader(body)).Decode(&handle); err != nil {
		return nil, err
	}
	fmt.Printf("hidriveGetMetaResponse prsed: %v\n", handle)

	return &handle, nil
}

func (hd *hiDrive) Open(name string) (drive.Handle, error) {
	return nil, nil
}

func (hd *hiDrive) OpenFile(name string, account *auth.Account) (*drive.File, error) {

	h, err := hd.GetMeta(name)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetHandle(%v)", name)
	}

	var file = drive.File{
		// Handle: h,
		URL:  path.Join(hd.URL, h.Path),
		Name: h.Name,
		Size: int64(h.Size),
		// // Mode:          h.mode,
		Type: h.GuessMIME(),
		// // Permissions:   h.mode.String(),
		Owner:         &drive.User{},  //GetUserByID(uid),
		Group:         &drive.Group{}, //GetGroupByID(gid),
		Authorization: h.GetPermissions(account),
		Modified:      time.Unix(h.MTime, 0),
	}

	return &file, nil
}
func (h *hiDrive) Create(string) (drive.Handle, error) {
	return nil, nil
}
func (h *hiDrive) Mkdir(string) (drive.Handle, error) {
	return nil, nil
}

func (h *hiDrive) ListFiles(folder *drive.File, account *auth.Account) ([]drive.File, error) {

	dir, err := hidriveGetDir(strings.TrimPrefix(folder.URL, h.URL), h.Token.GetAccessToken())
	if err != nil {
		fmt.Println("errdir", err)
		return nil, err
	}

	entries := make([]drive.File, dir.NMembers)

	for index, Member := range dir.Members {
		fmt.Println("size:", Member.Size, dir.Path)
		// 	h := Handle.(*handle)
		// 	uid, gid := h.userAndGroupIDs()
		entries[index] = drive.File{
			// 		Handle:        h,
			URL:  filepath.Join(folder.URL, Member.Name),
			Name: Member.Name,
			Size: int64(Member.Size),
			// 		Mode:          h.mode,
			Type: Member.GuessMIME(),
			// 		Permissions:   h.mode.String(),
			Owner:         &drive.User{},  //GetUserByID(uid),
			Group:         &drive.Group{}, //GetGroupByID(gid),
			Authorization: Member.GetPermissions(account),
			Modified:      time.Unix(Member.MTime, 0),
		}
	}
	return entries, nil
}