package fs

import (
	"mime"

	"github.com/ihleven/cloud11-api/drive"
)

func init() {

	mime.AddExtensionType(".py", "text/python")
	mime.AddExtensionType(".go", "text/golang")
	mime.AddExtensionType(".json", "text/json")
	mime.AddExtensionType(".js", "text/javascript")
	mime.AddExtensionType(".ts", "text/typescript")
	mime.AddExtensionType(".dia", "text/diary")
	mime.AddExtensionType(".md", "text/markdown")
}

var dir = drive.Type{
	Filetype:  "D",
	Mediatype: "dir",
	Subtype:   "",
	MIME:      "dir",
	Charset:   "",
}

var file = drive.Type{
	Filetype:  "F",
	Mediatype: "file",
	Subtype:   "",
	MIME:      "file",
	Charset:   "",
}

func (fh handle) GuessMIME() drive.Type {

	if fh.IsDir() {
		return dir
	}

	// ext := path.Ext(fh.Name())

	// if mimestr := mime.TypeByExtension(ext); mimestr != "" {
	// 	m = types.NewMIME(mimestr)

	// }

	// if m.Value == "" {
	// 	// m, _ = f.h2nonMatchMIME261()
	// }
	// if strings.HasSuffix(m.Subtype, "charset=utf-8") {
	// 	m.Subtype = m.Subtype[:len(m.Subtype)-15]
	// 	m.Value = m.Value[:len(m.Value)-15]
	// }
	return file
}
