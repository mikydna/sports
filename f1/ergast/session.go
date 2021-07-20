package ergast

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

var (
	ErrInvalidErgastSessionFunc = func(x interface{}) error {
		return fmt.Errorf("bad ergast: %v", x)
	}
)

const (
	YearURLToken   = "{year}"
	SeriesURLToken = "{series}"
	DateFormat     = "2006-01-02"
)

// generated
type Sessions struct {
	MRData struct {
		Xmlns     string `json:"xmlns"`
		Series    string `json:"series"`
		URL       string `json:"url"`
		Limit     string `json:"limit"`
		Offset    string `json:"offset"`
		Total     string `json:"total"`
		RaceTable struct {
			Season string `json:"season"`
			Races  []struct {
				Season   string `json:"season"`
				Round    string `json:"round"`
				URL      string `json:"url"`
				RaceName string `json:"raceName"`
				Circuit  struct {
					CircuitID   string `json:"circuitId"`
					URL         string `json:"url"`
					CircuitName string `json:"circuitName"`
					Location    struct {
						Lat      string `json:"lat"`
						Long     string `json:"long"`
						Locality string `json:"locality"`
						Country  string `json:"country"`
					} `json:"Location"`
				} `json:"Circuit"`
				Date string `json:"date"`
				Time string `json:"time"`
			} `json:"Races"`
		} `json:"RaceTable"`
	} `json:"MRData"`
}

// https://ergast.com/api/{series}/{year}.json
func (c *Client) Sessions(ctx context.Context, series string, season int) (*Sessions, error) {
	ep := sessionsURL(c.Host, series, season).String()
	resp, err := c.http.Get(ep)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var sessions *Sessions
	if err := json.NewDecoder(resp.Body).Decode(&sessions); err != nil {
		return nil, ErrInvalidErgastSessionFunc(err)
	}

	if reported := sessions.MRData.RaceTable.Season; reported != fmt.Sprint(season) {
		return nil, ErrInvalidErgastSessionFunc("unexpected season")
	}

	return sessions, nil
}

func sessionsURL(host, series string, season int) *url.URL {
	str := fmt.Sprintf("%s/%s/%s/%s", host, "api", "{series}", "{year}.json")
	str = strings.NewReplacer(
		SeriesURLToken, series,
		YearURLToken, fmt.Sprint(season),
	).Replace(str)

	sessionURL, err := url.Parse(str)
	if err != nil {
		panic(err)
	}

	return sessionURL
}
