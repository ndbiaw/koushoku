package server

import (
	"fmt"
	"math"
	"net/url"
	"strings"
	"time"

	. "koushoku/config"
	"koushoku/services"
)

var helper = map[string]any{
	"baseURL": func() string {
		return Config.Meta.BaseURL
	},

	"dataBaseURL": func() string {
		return Config.Meta.DataBaseURL
	},

	"language": func() string {
		return Config.Meta.Language
	},

	"formatBytes": func(v any) string {
		var b int64
		switch v := v.(type) {
		case int:
			b = int64(v)
		case int64:
			b = v
		case uint:
			b = int64(v)
		case uint64:
			b = int64(v)
		}
		return services.FormatBytes(b)
	},

	"formatNumber": func(v any) string {
		var n int64
		switch v := v.(type) {
		case int:
			n = int64(v)
		case int16:
			n = int64(v)
		case int32:
			n = int64(v)
		case int64:
			n = v
		case uint:
			n = int64(v)
		case uint16:
			n = int64(v)
		case uint32:
			n = int64(v)
		case uint64:
			n = int64(v)
		}
		return services.FormatNumber(n)
	},

	"joinURL": func(base string, s ...string) string {
		return services.JoinURL(base, s...)
	},

	"setQuery": func(query url.Values, key string, value any) string {
		query.Set(key, fmt.Sprintf("%v", value))
		return fmt.Sprintf("?%s", query.Encode())
	},

	"makeSlice": func(v any) []int {
		var n int
		switch v := v.(type) {
		case int:
			n = v
		case int16:
			n = int(v)
		case int32:
			n = int(v)
		case int64:
			n = int(v)
		case uint:
			n = int(v)
		case uint32:
			n = int(v)
		case uint64:
			n = int(v)
		}
		return make([]int, n)
	},

	"createQuery": func(query url.Values, key string, value any) string {
		clone := make(url.Values)
		for k, v := range query {
			clone[k] = v
		}
		clone.Set(key, fmt.Sprintf("%v", value))
		return fmt.Sprintf("?%s", clone.Encode())
	},

	"includes": func(slice []string, s string) bool {
		for _, v := range slice {
			if strings.EqualFold(v, s) {
				return true
			}
		}
		return false
	},

	"add": func(a, b int) int {
		return a + b
	},

	"sub": func(a, b int) int {
		return a - b
	},

	"mul": func(a, b int) int {
		return a * b
	},

	"div": func(a, b int) int {
		return a / b
	},

	"mod": func(a, b int) int {
		return a % b
	},

	"inc": func(a int) int {
		return a + 1
	},

	"dec": func(a int) int {
		return a - 1
	},

	"abs":   math.Abs,
	"floor": math.Floor,
	"ceil":  math.Ceil,
	"min":   math.Min,
	"max":   math.Max,

	"lowerCase":  strings.ToLower,
	"upperCase":  strings.ToUpper,
	"titleCase":  strings.Title,
	"trim":       strings.Trim,
	"trimLeft":   strings.TrimLeft,
	"trimRight":  strings.TrimRight,
	"trimSpace":  strings.TrimSpace,
	"trimPrefix": strings.TrimPrefix,
	"trimSuffix": strings.TrimSuffix,
	"hasPrefix":  strings.HasPrefix,
	"hasSuffix":  strings.HasSuffix,
	"contains":   strings.Contains,
	"replace":    strings.Replace,

	"formatTime": func(t time.Time, format string) string {
		return t.UTC().Format(format)
	},

	"formatUnix": func(n int64, format string) string {
		return time.Unix(n, 0).UTC().Format(format)
	},

	"currentTime": func() time.Time {
		return time.Now().UTC()
	},

	"currentUnix": func() int64 {
		return time.Now().UTC().Unix()
	},

	"currentYear": func() int {
		return time.Now().UTC().Year()
	},

	"currentMonth": func() int {
		return int(time.Now().UTC().Month())
	},

	"currentMonthString": func() string {
		return time.Now().UTC().Month().String()
	},

	"currentDay": func() int {
		return time.Now().UTC().Day()
	},

	"currentDayString": func() string {
		return time.Now().UTC().Weekday().String()
	},

	"currentHour": func() int {
		return time.Now().UTC().Hour()
	},

	"currentMinute": func() int {
		return time.Now().UTC().Minute()
	},

	"currentSecond": func() int {
		return time.Now().UTC().Second()
	},
}
