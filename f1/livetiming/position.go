package livetiming

import (
	"bytes"
	"encoding/json"
	"sort"
	"time"
)

// generated
type Position struct {
	Timestamp time.Time `json:"Timestamp"`
	Entries   map[uint]struct {
		Status string `json:"Status"`
		X      int    `json:"X"`
		Y      int    `json:"Y"`
		Z      int    `json:"Z"`
	} `json:"Entries"`
}

type PositionFrame struct {
	Offset   time.Duration `json:"offset"`
	Position []*Position   `json:"position"`
}

func (pf *PositionFrame) Unmarshal(b []byte) error {
	t, datab, err := DecodeLine(b, true)
	if err != nil {
		return err
	}

	var positions struct {
		Position []*Position `json:"Position"`
	}
	if err := json.Unmarshal(datab, &positions); err != nil {
		return err
	}

	*pf = PositionFrame{
		Offset:   t,
		Position: positions.Position,
	}

	return nil
}

type PositionFrames []*PositionFrame

func (pfs *PositionFrames) Unmarshal(b []byte) error {
	if pfs == nil || len(*pfs) != 0 {
		panic("frames are not empty")
	}

	lines := bytes.Split(bytes.TrimSpace(b), LineDelimiter)
	for _, line := range lines {
		position := new(PositionFrame)
		if err := position.Unmarshal(line); err != nil {
			return err
		}

		*pfs = append(*pfs, position) // nolint
	}

	// ensure frames are sorted by offset
	sort.Slice(*pfs, func(i, j int) bool {
		return (*pfs)[i].Offset.Nanoseconds() < (*pfs)[j].Offset.Nanoseconds()
	})

	return nil
}
