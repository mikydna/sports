package livetiming

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"io/ioutil"
	"time"

	"github.com/spkg/bom"
)

var (
	DateFormat        = "2006-01-02"
	SessionTimeFormat = "15:04:05.000"
	SessionTimeOffset = time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)
	LineDelimiter     = []byte("\n")
)

/** from FastF1
	{
		"archive_status": "ArchiveStatus.json",
		"audio_streams": "AudioStreams.jsonStream",
		"car_data": "CarData.z.jsonStream",
		"championship_predicion": "ChampionshipPrediction.jsonStream",
		"content_streams": "ContentStreams.jsonStream",
		"driver_list": "DriverList.jsonStream",
		"extrapolated_clock": "ExtrapolatedClock.jsonStream",
		"heartbeat": "Heartbeat.jsonStream",
		"lap_count": "LapCount.jsonStream",
		"position": "Position.z.jsonStream",
		"race_control_messages": "RaceControlMessages.json",
		"session_info": "SessionInfo.json",
		"session_status": "SessionStatus.jsonStream",
		"team_radio": "TeamRadio.jsonStream",
		"timing_app_data": "TimingAppData.jsonStream",
		"timing_data": "TimingData.jsonStream",
		"timing_stats": "TimingStats.jsonStream",
		"track_status": "TrackStatus.jsonStream",
		"weather_data": "WeatherData.jsonStream"
	}
**/

//go:generate stringer -type=File -linecomment=true -output livetiming_str.go
type File uint

const (
	FileSessionInfo   File = iota + 1 // SessionInfo.json
	FileDriverList                    // DriverList.jsonStream
	FilePosition                      // Position.z.jsonStream
	FileCarData                       // CarData.z.jsonStream
	FileTimingStats                   // TimingStats.jsonStream
	FileLapCount                      // LapCount.jsonStream
	FileTimingData                    // TimingData.jsonStream
	FileTimingAppData                 // TimingAppData.jsonStream
	FileTeamRadio                     // TeamRadio.jsonStream
)

var AllFiles []File = []File{
	FileSessionInfo, // basic race info
	FileDriverList,  // driver ids
	FilePosition,    // position (xyz)
	FileCarData,     // channels
	FileTimingStats,
	FileTimingAppData,
	FileTeamRadio,
	// FileLapCount,
	// FileTimingData,
}

func DecodeLine(b []byte, compressed bool) (time.Duration, []byte, error) {
	// livetimingdata is utf-bom; strip bom data
	cleanB := bom.Clean(b)

	// time (fixed size: 12b)
	t, err := time.Parse(SessionTimeFormat, string(cleanB[:12]))
	if err != nil {
		// panic(err)
		return -1, nil, err
	}

	offset := t.Sub(SessionTimeOffset)

	// data
	datab := bytes.Trim(cleanB[12:], `"`)
	if compressed {
		_, err := base64.StdEncoding.Decode(datab, datab)
		if err != nil {
			return -1, nil, err
		}

		r := flate.NewReader(bytes.NewBuffer(datab))
		datab, err = ioutil.ReadAll(r)
		if err != nil {
			return -1, nil, err
		}

	}

	return offset, datab, nil
}
