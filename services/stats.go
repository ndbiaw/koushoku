package services

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"

	. "koushoku/config"
)

type Stats struct {
	ArchiveCount  int64
	ArtistCount   int64
	CircleCount   int64
	MagazineCount int64
	ParodyCount   int64
	TagCount      int64

	PageCount        int64
	AveragePageCount int64
	Size             int64
	AverageSize      int64

	Analytics *Analytics

	sync.RWMutex
	sync.Once
}

var stats Stats

func AnalyzeStats() (err error) {
	log.Println("Analyzing stats...")
	defer func() {
		if err != nil {
			log.Println("AnalyzeStats returned an error:", err)
		}
	}()

	stats.ArchiveCount, err = GetArchiveCount()
	if err != nil {
		return
	}

	stats.Size, stats.PageCount, err = GetArchiveStats()
	if err != nil {
		return
	}

	stats.ArtistCount, err = GetArtistCount()
	if err != nil {
		return
	}

	stats.CircleCount, err = GetCircleCount()
	if err != nil {
		return
	}
	stats.MagazineCount, err = GetMagazineCount()
	if err != nil {
		return
	}

	stats.ParodyCount, err = GetParodyCount()
	if err != nil {
		return
	}

	stats.TagCount, err = GetTagCount()
	if err != nil {
		return
	}

	if stats.ArchiveCount > 0 {
		if stats.PageCount > 0 {
			stats.AveragePageCount = int64(math.Round(float64(stats.PageCount) / float64(stats.ArchiveCount)))
		}
		if stats.Size > 0 {
			stats.AverageSize = int64(math.Round(float64(stats.Size) / float64(stats.ArchiveCount)))
		}
	}

	if Config.Mode == "production" {
		err = fetchAnalytics()
		if err != nil {
			return
		}

		stats.Do(func() {
			go func() {
				for {
					time.Sleep(30 * time.Minute)
					log.Println("Refreshing analytics...")

					stats.Lock()
					if err = fetchAnalytics(); err != nil {
						log.Println("Failed to refresh analytics", err)
					}
					stats.Unlock()
				}
			}()
		})
	}

	return
}

func GetStats() *Stats {
	stats.RLock()
	defer stats.RUnlock()
	return &stats
}

type Analytic struct {
	Date           string
	Bytes          int64
	CachedBytes    int64
	Requests       int64
	CachedRequests int64
}

type Analytics struct {
	Analytic
	Entries     []*Analytic
	LastUpdated time.Time
}

type GraphQL struct {
	OperationName string         `json:"operationName,omitempty"`
	Query         string         `json:"query"`
	Variables     map[string]any `json:"variables"`
}

func (g *GraphQL) Marshal() *strings.Reader {
	buf, _ := json.Marshal(g)
	return strings.NewReader(string(buf))
}

func fetchAnalytics() error {
	payload := &GraphQL{
		OperationName: "GetZoneAnalytics",
		Query: `query GetZoneAnalytics($zoneTag: string, $since: string, $until: string) {
			viewer {
				zones(filter: {zoneTag: $zoneTag}) {
					totals: httpRequests1dGroups(limit: 10000, filter: {date_geq: $since, date_lt: $until}) {
						sum {
							bytes
							cachedBytes
							requests
							cachedRequests
						}
					}
					zones: httpRequests1dGroups(orderBy: [date_ASC], limit: 10000, filter: {date_geq: $since, date_lt: $until}) {
						dimensions {
							timeslot: date
						}
						sum {
							bytes
							cachedBytes
							requests
							cachedRequests
						}
					}
				}
			}
		}`,
		Variables: map[string]any{
			"zoneTag": Config.Cloudflare.ZoneTag,
			"since":   fmt.Sprintf("%d-01-01", time.Now().Year()),
			"until":   time.Now().AddDate(0, 1, 0).Format("2006-01-02"),
		},
	}

	u := "https://api.cloudflare.com/client/v4/graphql"
	req, err := http.NewRequest("POST", u, payload.Marshal())
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Email", Config.Cloudflare.Email)
	req.Header.Set("X-Auth-Key", Config.Cloudflare.ApiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body := &struct {
		Data struct {
			Viewer struct {
				Zones []struct {
					Totals []struct {
						Sum struct {
							Bytes          int64 `json:"bytes"`
							CachedBytes    int64 `json:"cachedBytes"`
							Requests       int64 `json:"requests"`
							CachedRequests int64 `json:"cachedRequests"`
						} `json:"sum"`
					} `json:"totals"`
					Zones []struct {
						Dimensions struct {
							Timeslot string `json:"timeslot"`
						} `json:"dimensions"`
						Sum struct {
							Bytes          int64 `json:"bytes"`
							CachedBytes    int64 `json:"cachedBytes"`
							Requests       int64 `json:"requests"`
							CachedRequests int64 `json:"cachedRequests"`
						} `json:"sum"`
					} `json:"zones"`
				} `json:"zones"`
			} `json:"viewer"`
		} `json:"data"`
	}{}

	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return err
	}

	if stats.Analytics == nil {
		stats.Analytics = &Analytics{}
	}

	analytics := stats.Analytics
	prevBytes := analytics.Bytes
	prevCachedBytes := analytics.CachedBytes
	prevRequests := analytics.Requests
	prevCachedRequests := analytics.CachedRequests

	if len(body.Data.Viewer.Zones) > 0 {
		if len(body.Data.Viewer.Zones[0].Totals) > 0 {
			analytics.Bytes = body.Data.Viewer.Zones[0].Totals[0].Sum.Bytes
			analytics.CachedBytes = body.Data.Viewer.Zones[0].Totals[0].Sum.CachedBytes
			analytics.Requests = body.Data.Viewer.Zones[0].Totals[0].Sum.Requests
			analytics.CachedRequests = body.Data.Viewer.Zones[0].Totals[0].Sum.CachedRequests
		}

		analytics.Entries = []*Analytic{}
		for _, zone := range body.Data.Viewer.Zones {
			for _, entry := range zone.Zones {
				analytics.Entries = append(analytics.Entries, &Analytic{
					Date:           entry.Dimensions.Timeslot,
					Bytes:          entry.Sum.Bytes,
					CachedBytes:    entry.Sum.CachedBytes,
					Requests:       entry.Sum.Requests,
					CachedRequests: entry.Sum.CachedRequests,
				})
			}
		}
	}

	if prevBytes != analytics.Bytes || prevCachedBytes != analytics.CachedBytes ||
		prevRequests != analytics.Requests || prevCachedRequests != analytics.CachedRequests {
		analytics.LastUpdated = time.Now()
	}
	return nil
}
