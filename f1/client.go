package f1

import (
	"io"
	"net/http"
	"os"

	"github.com/mikydna/sports/f1/ergast"
	"github.com/mikydna/sports/f1/livetiming"
)

type Client struct {
	*DownloadService
	*RepoService
	w io.Writer
}

func NewClient(c *http.Client, w io.Writer, repo, workspace string) *Client {
	eg := ergast.NewClient(c, ergast.DefaultHost)
	lt := livetiming.NewClient(c, livetiming.DefaultHost)
	return &Client{
		DownloadService: &DownloadService{repo, eg, lt, w},
		RepoService:     &RepoService{repo, workspace},
		w:               w,
	}
}

func DefaultClient(repo, workspace string) *Client {
	return NewClient(http.DefaultClient, os.Stdout, repo, workspace)
}
