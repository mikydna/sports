package f1

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	ErrExportUnknownData       = errors.New("export: unknown data type")
	ErrExportUnsupportedFormat = errors.New("export: unsupported format")
)

//go:generate stringer -type=ExportFormat -linecomment=true -output export_str.go
type ExportFormat uint8

const (
	ExportFormatUnknown ExportFormat = iota // unknown
	ExportFormatJSON                        // jsonl
	ExportFormatProto                       // pb
	ExportFormatGob                         // gob
)

type ExportService struct {
	workspace string
}

func NewExportService(workspace string) *ExportService {
	return &ExportService{workspace}
}

func (s *ExportService) Export(ctx context.Context, dst string, data interface{}, format ExportFormat) (string, error) {
	dir := filepath.Clean(filepath.Join(s.workspace, dst+"/")) // ensure its a dir
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	switch val := data.(type) {
	case PositionFrames:
		fp := filepath.Join(dir, fmt.Sprintf("%s.%v", "Position", format))

		switch format {
		case ExportFormatJSON:
			if err := exportJSON(fp, data); err != nil {
				return "", err
			}
		case ExportFormatGob:
			if err := exportGob(fp, data); err != nil {
				return "", err
			}
		case ExportFormatProto:
			b, err := val.ProtoBytes()
			if err != nil {
				return "", err
			}

			if err := os.WriteFile(fp, b, 0644); err != nil {
				return "", err
			}

			return fp, nil

		default:
			return "", ErrExportUnsupportedFormat
		}

	default:
		return "", ErrExportUnknownData
	}

	return "", nil
}

func exportJSON(fp string, data interface{}) error {
	f, err := os.Open(fp)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(data); err != nil {
		return err
	}

	return nil
}

func exportGob(fp string, data interface{}) error {
	f, err := os.Open(fp)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := gob.NewEncoder(f).Encode(data); err != nil {
		return err
	}

	return nil
}
