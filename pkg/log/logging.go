package log

import (
	"net/url"
	"time"

	"github.com/fatih/color"
)

func Info(format string, args ...interface{}) {
	color.Cyan(" ---- "+format+"\n", args...)
	color.Cyan(" asdf")
}

func Access(url url.URL, statusCode int, duration time.Duration, format string, args ...interface{}) {
	color.Cyan(" * %s?%s   Status: %v, took: %v => "+format+"\n", url.Path, url.RawQuery, statusCode, duration)
	color.Cyan(" asdf")

}