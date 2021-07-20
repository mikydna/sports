package f1

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"

	"github.com/mikydna/sports/f1/livetiming"
)

type RepoService struct {
	Repo      string
	Workspace string
}

func (s *RepoService) List(ctx context.Context, w io.Writer) error {
	sessionFiles, err := searchRepo(s.Repo, livetiming.FileSessionInfo)
	if err != nil {
		return err
	}

	for _, curr := range sessionFiles {
		rel, err := filepath.Rel(s.Repo, curr)
		if err != nil {
			return err
		}

		fmt.Fprintln(w, filepath.Dir(rel))
	}

	return nil
}

func searchRepo(repo string, target livetiming.File) ([]string, error) {
	result := []string{}
	if err := filepath.WalkDir(repo, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() && d.Name() == ".git" {
			return filepath.SkipDir
		}
		if d.Name() == target.String() {
			result = append(result, path)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}
