package hidrive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// https://my.hidrive.com/client/authorize?client_id=b4436f1157043c2bf8db540c9375d4ed&response_type=code&scope=admin,rw

// => o67hRuKoho4KSunvhSmj

// {
// 	"refresh_token":"rt-cmgrqgjgodc1xomgfsdc",
// 	"expires_in":3600,
// 	"userid":"59995203.2308.9433",
// 	"access_token":"z9N59zQpyWt4JI0h4fRs",
// 	"alias":"ihleven",
// 	"token_type":"Bearer"
// }

func getAccessToken(code string) (string, error) {

	tokenbytes, err := ioutil.ReadFile("./token")
	if err == nil {
		fmt.Println("read:", string(tokenbytes))
		return string(tokenbytes), nil
	}

	resp, err := hidriveOAuth2Token(code)
	if err != nil {
		fmt.Println("oauth error:", err)
	}

	err = ioutil.WriteFile("./token", []byte(resp.AccessToken), 0644)
	if err != nil {
		return resp.AccessToken, err
	}
	return resp.AccessToken, nil
}

type OAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
}

func hidriveOAuth2Token(code string) (*OAuthAccessResponse, error) {

	formData := url.Values{
		"client_id":     {"b4436f1157043c2bf8db540c9375d4ed"},
		"client_secret": {"8c5453a7264e4200ab80206658987dd8"},
		"grant_type":    {"authorization_code"},
		"code":          {code},
	}

	res, err := http.PostForm("https://my.hidrive.com/oauth2/token", formData)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var response OAuthAccessResponse
	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}

// curl -X POST --data // https://my.hidrive.com/oauth2/token

type DirResponse struct {
	Path     string     `json:"path"`
	Type     string     `json:"type"`
	Size     uint64     `json:"size"`
	Readable bool       `json:"readable"`
	Writable bool       `json:"writable"`
	CTime    int64      `json:"ctime"`
	MTime    int64      `json:"mtime"`
	HasDirs  bool       `json:"has_dirs"`
	NMembers int        `json:"nmembers"`
	Members  []hiHandle `json:"members"`
}

type Image struct {
	Width  int  `json:"width"`
	Height int  `json:"height"`
	Exif   Exif `json:"exif"`
}

type Exif struct {
	DateTimeOriginal string  `json:"DateTimeOriginal"`
	Make             string  `json:"Make"`
	Model            string  `json:"Model"`
	ImageWidth       int     `json:"ImageWidth"`
	ImageHeight      int     `json:"ImageHeight"`
	ExifImageWidth   int     `json:"ExifImageWidth"`
	ExifImageHeight  int     `json:"ExifImageHeight"`
	XResolution      float64 `json:"XResolution"`
	YResolution      float64 `json:"YResolution"`
	ResolutionUnit   int     `json:"ResolutionUnit"`
	BitsPerSample    int     `json:"BitsPerSample"`
	Aperture         float64 `json:"Aperture"`
	ExposureTime     float64 `json:"ExposureTime"`
	ISO              int     `json:"ISO"`
	FocalLength      float64 `json:"FocalLength"`
	Orientation      float64 `json:"Orientation"`
	GPSLatitude      float64 `json:"GPSLatitude"`
	GPSLongitude     float64 `json:"GPSLongitude"`
	GPSAltitude      float64 `json:"GPSAltitude"`
}

func hidriveGetDir(path string, bearer string) (*DirResponse, error) {

	members := "members.id,members.mime_type,members.mtime,members.name,members.readable,members.writable,members.type,members.nmembers,members.path,members.size"
	queryParams := url.Values{
		"path":    {path},
		"members": {"all"},
		"fields":  {"ctime,has_dirs,id,mtime,readable,size,type,writable,nmembers,path," + members},
	}
	req, _ := http.NewRequest("GET", "https://api.hidrive.strato.com/2.1/dir", nil)
	req.URL.RawQuery = queryParams.Encode()
	req.Header.Set("Authorization", "Bearer "+bearer)
	client := &http.Client{}
	res, err := client.Do(req)
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
	//fmt.Println("hidriveGetDirResponse:", string(body), res.Status, res.StatusCode, "adf")

	var response DirResponse
	if err = json.NewDecoder(bytes.NewReader(body)).Decode(&response); err != nil {
		return nil, err
	}
	fmt.Printf("dir response: %v\n", response)

	return &response, nil
}

type hidriveMetaResponse struct {
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
}

type hidriveError struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

func (e hidriveError) Error() string {
	return e.Message
}

func NewHiDriveError(txt io.ReadCloser, code int, status string) error {
	var hidriveError hidriveError
	if err := json.NewDecoder(txt).Decode(&hidriveError); err == nil {
		return &hidriveError
	}

	body, err := ioutil.ReadAll(txt)
	if err != nil {
		return err
	}
	hidriveError.Code = code
	hidriveError.Message = status
	fmt.Println("body:", body, code, status)
	return &hidriveError
}
func hidriveGetFile(path string, bearer string, w http.ResponseWriter) error {
	// 200 OK, 206 Partial content
	queryParams := url.Values{
		"path": {path},
	}

	request, _ := http.NewRequest("GET", "https://api.hidrive.strato.com/2.1/file", nil)
	request.URL.RawQuery = queryParams.Encode()
	request.Header.Set("Authorization", "Bearer "+bearer)

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return &hidriveError{res.StatusCode, res.Status}
	}
	if res.StatusCode != 200 {
		fmt.Println(res.Body)
		return NewHiDriveError(res.Body, res.StatusCode, res.Status)
	}
	if _, err := io.Copy(w, res.Body); err != nil {
		fmt.Println("copy:", err)
		return err
	}
	fmt.Println("file created")
	return nil
}
