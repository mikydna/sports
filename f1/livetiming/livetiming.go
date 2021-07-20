package livetiming

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	DateFormat = "2006-01-02"
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
	FileSessionInfo File = iota + 1 // SessionInfo.json
	FileDriverList                  // DriverList.jsonStream
	FilePosition                    // Position.z.jsonStream
	FileCarData                     // CarData.z.jsonStream
	FileTimingStats                 // TimingStats.jsonStream
	FileLapCount                    // LapCount.jsonStream
	FileTimingData                  // TimingData.jsonStream
	FileTeamRadio                   // TeamRadio.jsonStream
)

var AllFiles []File = []File{
	FileSessionInfo, // basic race info
	FileDriverList,  // driver ids
	FilePosition,    // position (xyz)
	FileCarData,     // channels
	FileTimingStats,
	FileTeamRadio,
	// FileLapCount,
	// FileTimingData,
}

func normPath(s string) string {
	fn := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	x, _, e := transform.String(fn, s)
	if e != nil {
		panic(e)
	}
	return strings.Replace(x, " ", "_", -1)
}
