package livetiming

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
)

const (
	DefaultHost = "https://livetiming.formula1.com"
)

var (
	DefaultDownloadOptions = &DownloadOptions{
		SkipIfExists: true,
	}
)

type Client struct {
	http *http.Client
	Host string
}

type DownloadProgress struct {
	Src     string
	Dst     string
	Skipped bool
	Err     error
	Pct     float64
}

type DownloadOptions struct {
	SkipIfExists bool
	Progress     chan<- *DownloadProgress
}

func NewClient(c *http.Client, host string) *Client {
	return &Client{c, host}
}

func DefaultClient() *Client {
	return NewClient(http.DefaultClient, DefaultHost)
}

func (c *Client) DownloadFiles(ctx context.Context, dst string, date time.Time, name string, files []File, opt *DownloadOptions) error {
	if opt == nil {
		opt = DefaultDownloadOptions
	}

	uri := raceURI(date, name)
	downloadURL, _ := url.Parse(DefaultHost)
	downloadURL.Path = strings.Join([]string{"static", uri}, "/")

	var errs error
	for i, file := range files {
		src := strings.Join([]string{downloadURL.String(), file.String()}, "/")
		dst := filepath.Join(dst, uri, file.String())
		pct := float64(i+1) / float64(len(files))

		// skip if exists
		if opt.SkipIfExists && fileExists(dst) {
			progressDownload(opt.Progress, src, dst, true, nil, pct)
			continue
		}

		err := raceDownload(c.http, src, dst)
		if err != nil {
			errs = multierror.Append(errs, err)
		}

		progressDownload(opt.Progress, src, dst, false, err, pct)
	}

	return nil
}

func raceURI(raceDate time.Time, raceName string) string {
	seasonStr := fmt.Sprint(raceDate.Year())
	dateStr := raceDate.Format(DateFormat)
	weekendStr := fmt.Sprintf("%s_%s", dateStr, raceName)
	sessionStr := fmt.Sprintf("%s_%s", dateStr, "Race")
	return normPath(strings.Join([]string{seasonStr, weekendStr, sessionStr}, "/"))
}

func raceDownload(c *http.Client, src, dst string) error {
	// fetch data
	resp, err := c.Get(src)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// write file
	if err := os.MkdirAll(filepath.Dir(dst), 0644); err != nil {
		return err
	}

	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return err
	}

	return nil
}

func progressDownload(c chan<- *DownloadProgress, src, dst string, skipped bool, err error, pct float64) {
	if c == nil {
		return
	}

	select {
	case c <- &DownloadProgress{src, dst, skipped, err, pct}:
	default:
	}
}
