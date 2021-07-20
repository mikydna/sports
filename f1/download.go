package f1

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/mikydna/sports/f1/ergast"
	"github.com/mikydna/sports/f1/livetiming"
	"golang.org/x/sync/errgroup"
)

var (
	ErrF1DownloadConfigFunc = func(x interface{}) error {
		return fmt.Errorf("bad download config: %v", x)
	}
)

type DownloadConfig struct {
	Seasons []int
}

func checkConfig(cfg *DownloadConfig) error {
	if cfg == nil {
		return ErrF1DownloadConfigFunc("nil config")
	}
	// if !dirExists(cfg.Dest) {
	// 	return ErrF1DownloadConfigFunc("invalid dir")
	// }
	return nil
}

type DownloadService struct {
	repo string
	eg   *ergast.Client
	lt   *livetiming.Client
	w    io.Writer
}

func (s *DownloadService) Download(ctx context.Context, cfg *DownloadConfig) error {
	if err := checkConfig(cfg); err != nil {
		return err
	}

	var errs error
	grp, grpCtx := errgroup.WithContext(ctx)
	for _, season := range cfg.Seasons {
		sessions, err := s.eg.Sessions(ctx, "f1", season)
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}

		for _, race := range sessions.MRData.RaceTable.Races {
			raceName := race.RaceName
			raceDate, err := time.Parse(ergast.DateFormat, race.Date)
			if err != nil {
				errs = multierror.Append(errs, err)
				continue
			}

			// skip if race date is in the future
			if raceDate.After(time.Now()) {
				continue
			}

			// run download in thread
			grp.Go(func() error {
				downloadOpts := &livetiming.DownloadOptions{
					SkipIfExists: false,
					Progress:     nil,
				}

				if err := s.lt.DownloadFiles(grpCtx, s.repo, raceDate, raceName, livetiming.AllFiles, downloadOpts); err != nil {
					errs = multierror.Append(errs, err)
					return nil
				}

				log.Println("OK", raceDate, raceName)

				return nil
			})
		}
	}

	_ = grp.Wait()

	return errs
}
