package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	. "koushoku/config"
	"koushoku/modext"

	"github.com/gosimple/slug"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type IndexMap struct {
	Cache map[string]bool
	sync.RWMutex
}

func (m *IndexMap) Add(k string, v bool) {
	m.Lock()
	m.Cache[k] = v
	m.Unlock()
}

func (m *IndexMap) AddIfNotExist(k string, v bool) {
	if !m.Has(k) {
		m.Add(k, v)
	}
}

func (m *IndexMap) Clear() {
	m.Lock()
	m.Cache = make(map[string]bool)
	m.Unlock()
}

func (m *IndexMap) Get(key string) (bool, bool) {
	m.RLock()
	defer m.RUnlock()
	v, ok := m.Cache[key]
	return v, ok
}

func (m *IndexMap) Has(key string) bool {
	m.RLock()
	defer m.RUnlock()
	_, ok := m.Cache[key]
	return ok
}

func (m *IndexMap) Remove(key string) {
	m.Lock()
	delete(m.Cache, key)
	m.Unlock()
}

type Pagination struct {
	CurrentPage int
	Pages       []int
	TotalPages  int
}

const maxPages = 10

func CreatePagination(currentPage, totalPages int) *Pagination {
	if currentPage < 1 {
		currentPage = 1
	} else if currentPage > totalPages {
		currentPage = totalPages
	}

	pagination := &Pagination{
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	}

	var first, last int
	if totalPages <= maxPages {
		first = 1
		last = totalPages
	} else {
		min := int(math.Floor(float64(maxPages) / 2))
		max := int(math.Ceil(float64(maxPages)/2)) - 1
		if currentPage <= min {
			first = 1
			last = maxPages
		} else if currentPage+max >= totalPages {
			first = totalPages - maxPages + 1
			last = totalPages
		} else {
			first = currentPage - min
			last = currentPage + max
		}
	}

	pagination.Pages = make([]int, last-first+1)
	for i := 0; i < last+1-first; i++ {
		pagination.Pages[i] = first + i
	}

	return pagination
}

func FormatBytes(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

func FormatNumber(n int64) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d", n)
}

func FileName(path string) string {
	return strings.TrimRight(filepath.Base(path), filepath.Ext(path))
}

var pageNumRgx = regexp.MustCompile("[0-9]+")

func GetPageNum(fileName string) int {
	fileName = strings.TrimLeft(fileName, "0")
	n, _ := strconv.Atoi(pageNumRgx.FindString(fileName))
	return n
}

var imgFormatRgx = regexp.MustCompile(`(?i)(gif|jpe?g|tiff?|png|webp|bmp)$`)

func IsImage(path string) bool {
	return imgFormatRgx.MatchString(path)
}

func JoinURL(base string, paths ...string) string {
	u, _ := url.Parse(base)
	for _, path := range paths {
		u.Path = filepath.Join(u.Path, strings.TrimLeft(strings.TrimRight(path, "/"), "/"))
	}
	return u.String()
}

func JoinOR(strs ...string) string {
	return fmt.Sprintf("(%s)", strings.Join(strs, " OR "))
}

func makeCacheKey(v any) string {
	buf, _ := json.Marshal(v)
	return string(buf)
}

// Min returns the smaller of x or y.
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max returns the larger of x or y.
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

var slugCache struct {
	Map map[string]string
	sync.RWMutex
	sync.Once
}

func init() {
	slug.CustomSub = make(map[string]string)
	slug.CustomSub["❤"] = ""
	slug.CustomSub["☆"] = "-"
	slug.CustomSub["&"] = ""
	slug.CustomSub["♀"] = "bjb"
	slug.CustomSub["_"] = "-"
}

func Slugify(s string) string {
	slugCache.Once.Do(func() {
		slugCache.Map = make(map[string]string)
	})

	if len(s) == 0 {
		return ""
	}
	s = strings.ToLower(s)

	slugCache.RLock()
	v, ok := slugCache.Map[s]
	slugCache.RUnlock()

	if !ok {
		v = slug.Make(s)

		slugCache.Lock()
		slugCache.Map[s] = v
		slugCache.Unlock()
	}
	return v
}

func SlugifyStrings(strs []string) []string {
	for i, s := range strs {
		strs[i] = Slugify(s)
	}
	sort.Strings(strs)
	return strs
}

func Pluralize(str string) string {
	if strings.HasSuffix(str, "ss") {
		return str + "es"
	} else if strings.HasSuffix(str, "y") {
		return strings.TrimSuffix(str, "y") + "ies"
	}
	return str + "s"
}

type ResizeOptions struct {
	Width  int
	Height int
	PNG    bool
}

var resizer struct {
	Map   map[string]*sync.Mutex
	Queue chan bool
	sync.RWMutex
	sync.Once
}

func init() {
	resizer.Map = make(map[string]*sync.Mutex)
	resizer.Queue = make(chan bool, 10)
}

func ResizeImage(filePath, outputPath string, opts ResizeOptions) error {
	resizer.RLock()
	mu, ok := resizer.Map[outputPath]
	resizer.RUnlock()

	if !ok {
		mu = &sync.Mutex{}

		resizer.Lock()
		resizer.Map[outputPath] = mu
		resizer.Unlock()
	}

	mu.Lock()
	defer func() {
		mu.Unlock()

		resizer.Lock()
		delete(resizer.Map, outputPath)
		resizer.Unlock()
	}()

	if ok {
		return nil
	}

	resizer.Queue <- true
	defer func() {
		<-resizer.Queue
	}()

	args := []string{
		"-thumbnail", fmt.Sprintf("%dx%d^", opts.Width, opts.Height),
		"-gravity", "Center",
		"-extent", fmt.Sprintf("%dx%d", opts.Width, opts.Height),
		"-quality", "85",
		"-write", outputPath,
		filePath}

	var err error
	for i := 0; i < 3; i++ {
		if i == 1 {
			args = append(args, "-define", "colorspace:auto-grayscale=false")
		} else if i == 2 {
			args = append(args[:len(args)-3], fmt.Sprintf("PNG32:%s", filePath))
		}

		if _, err = RunCommand("mogrify", args...); err != nil && opts.PNG {
			continue
		}
		break
	}

	_, err = os.Stat(outputPath)
	return err
}

func RunCommand(path string, args ...string) (*bytes.Buffer, error) {
	cmd := exec.Command(path, args...)

	var buf bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &buf
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	if err := stderr.String(); len(err) > 0 {
		return nil, errors.New(err)
	}
	return &buf, nil
}

func GetArchiveSymlink(id int) (string, error) {
	symlink := filepath.Join(Config.Directories.Symlinks, strconv.Itoa(id))
	return os.Readlink(symlink)
}

func CreateArchiveSymlink(archive *modext.Archive) error {
	if archive == nil {
		return nil
	}

	if _, err := os.Stat(Config.Directories.Symlinks); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(Config.Directories.Symlinks, 0755); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	symlink := filepath.Join(Config.Directories.Symlinks, strconv.Itoa(int(archive.ID)))
	return os.Symlink(archive.Path, symlink)
}
