package config

import (
	_ "embed"

	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/google/uuid"
	"github.com/jessevdk/go-flags"
	"gopkg.in/ini.v1"
)

//go:embed config.ini
var buf []byte

var opts struct {
	Path     string `short:"c" description:"Path to config file"`
	Mode     string `short:"m" description:"App mode"`
	WebPort  int    `short:"p" description:"Web server port"`
	DataPort int    `short:"d" description:"Data server port"`
}

var Config struct {
	file *ini.File
	mu   sync.RWMutex

	Mode string

	Meta struct {
		BaseURL     string
		DataBaseURL string
		Title       string
		Description string
		Language    string
	}

	Database struct {
		Host    string
		Port    int
		Name    string
		User    string
		Passwd  string
		SSLMode string
	}

	Redis struct {
		Host   string
		Port   int
		DB     int
		Passwd string
	}

	Server struct {
		WebPort  int
		DataPort int
	}

	HTTP struct {
		Cookie string
		ApiKey string
	}

	Cloudflare struct {
		Email   string
		ApiKey  string
		ZoneTag string
	}

	Directories struct {
		Root string
		Data string

		Symlinks   string
		Templates  string
		Thumbnails string
	}

	Paths struct {
		Alias     string
		Blacklist string
		Metadata  string
	}
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	exec, err := os.Executable()
	if err == nil {
		exec, err = filepath.EvalSymlinks(exec)
	}
	if err != nil {
		log.Fatalln(err)
	}

	flags.NewParser(&opts, flags.PassDoubleDash).Parse()

	Config.Directories.Root = filepath.Dir(exec)
	Config.Directories.Symlinks = filepath.Join(Config.Directories.Root, "symlinks")
	if _, err := os.Stat(Config.Directories.Symlinks); os.IsNotExist(err) {
		log.Println("Creating symlinks directory...")
		if err := os.Mkdir(Config.Directories.Symlinks, 0755); err != nil {
			log.Fatalln(err)
		}
	}

	Config.Directories.Templates = filepath.Join(Config.Directories.Root, "templates")
	if _, err := os.Stat(Config.Directories.Templates); os.IsNotExist(err) {
		log.Println("Creating templates directory...")
		if err := os.Mkdir(Config.Directories.Templates, 0755); err != nil {
			log.Fatalln(err)
		}
	}

	Config.Directories.Thumbnails = filepath.Join(Config.Directories.Root, "thumbnails")
	if _, err := os.Stat(Config.Directories.Thumbnails); os.IsNotExist(err) {
		log.Println("Creating thumbnails directory...")
		if err := os.MkdirAll(Config.Directories.Thumbnails, 0755); err != nil {
			log.Fatalln(err)
		}
	}

	Config.Paths.Blacklist = filepath.Join(Config.Directories.Root, "blacklist.txt")
	if _, err := os.Stat(Config.Paths.Blacklist); os.IsNotExist(err) {
		log.Println("No blacklist file found, creating one...")
		if err := os.WriteFile(Config.Paths.Blacklist, []byte(""), 0755); err != nil {
			log.Fatalln(err)
		}
	}

	Config.Paths.Metadata = filepath.Join(Config.Directories.Root, "metadata.json")
	if _, err := os.Stat(Config.Paths.Metadata); os.IsNotExist(err) {
		log.Println("No metadata file found, creating one...")
		if err := os.WriteFile(Config.Paths.Metadata, []byte("{}"), 0755); err != nil {
			log.Fatalln(err)
		}
	}

	Config.Paths.Alias = filepath.Join(Config.Directories.Root, "alias.txt")
	if _, err := os.Stat(Config.Paths.Alias); os.IsNotExist(err) {
		log.Println("No alias file found, creating one...")
		if err := os.WriteFile(Config.Paths.Alias, []byte("{}"), 0755); err != nil {
			log.Fatalln(err)
		}
	}

	if len(opts.Path) == 0 {
		opts.Path = filepath.Join(Config.Directories.Root, "config.ini")
	}

	if _, err = os.Stat(opts.Path); os.IsNotExist(err) {
		log.Println("No config file found, creating one...")
		if err := os.MkdirAll(filepath.Dir(opts.Path), 0755); err != nil {
			log.Fatalln(err)
		} else if err := os.WriteFile(opts.Path, buf, 0755); err != nil {
			log.Fatalln(err)
		}
	}

	var file *ini.File
	if file, err = ini.Load(opts.Path); err != nil {
		log.Fatalln(err)
	}

	Config.file = file
	Config.Mode = file.Section("").Key("mode").MustString("production")

	Config.Meta.BaseURL = file.Section("meta").Key("base_url").MustString("http://localhost:42073")
	Config.Meta.DataBaseURL = file.Section("meta").Key("data_base_url").MustString("http://localhost:42075")
	Config.Meta.Title = file.Section("meta").Key("title").MustString("Koushoku")
	Config.Meta.Description = file.Section("meta").Key("description").String()
	Config.Meta.Language = file.Section("meta").Key("language").MustString("en-US")

	Config.Database.Host = file.Section("database").Key("host").MustString("localhost")
	Config.Database.Port = file.Section("database").Key("port").MustInt(5432)
	Config.Database.Name = file.Section("database").Key("name").MustString("koushoku")
	Config.Database.User = file.Section("database").Key("user").MustString("koushoku")
	Config.Database.Passwd = file.Section("database").Key("passwd").MustString("koushoku")
	Config.Database.SSLMode = file.Section("database").Key("ssl_mode").MustString("disable")

	Config.Redis.Host = file.Section("redis").Key("host").MustString("localhost")
	Config.Redis.Port = file.Section("redis").Key("port").MustInt(6379)
	Config.Redis.DB = file.Section("redis").Key("db").MustInt(0)
	Config.Redis.Passwd = file.Section("redis").Key("passwd").String()

	Config.Server.WebPort = 42073
	Config.Server.DataPort = 42075
	Config.HTTP.Cookie = file.Section("http").Key("cookie").String()
	Config.HTTP.ApiKey = file.Section("http").Key("api_key").MustString(uuid.NewString())

	Config.Cloudflare.Email = file.Section("cloudflare").Key("email").MustString("")
	Config.Cloudflare.ApiKey = file.Section("cloudflare").Key("api_key").MustString("")
	Config.Cloudflare.ZoneTag = file.Section("cloudflare").Key("zone_tag").MustString("")

	Config.Directories.Data = file.Section("directories").Key("data").
		MustString(filepath.Join(Config.Directories.Root, "data"))
	if _, err := os.Stat(Config.Directories.Data); os.IsNotExist(err) {
		log.Println("Creating data directory...")
		if err := os.Mkdir(Config.Directories.Data, 0755); err != nil {
			log.Fatalln(err)
		}
	}

	Save()

	if len(opts.Mode) > 0 {
		Config.Mode = opts.Mode
	}

	if opts.WebPort > 0 {
		Config.Server.WebPort = opts.WebPort
	}

	if opts.DataPort > 0 {
		Config.Server.DataPort = opts.DataPort
	}
}

func Save() error {
	Config.mu.Lock()
	defer Config.mu.Unlock()

	Config.file.Section("").Key("mode").SetValue(Config.Mode)

	Config.file.Section("meta").Key("base_url").SetValue(Config.Meta.BaseURL)
	Config.file.Section("meta").Key("data_base_url").SetValue(Config.Meta.DataBaseURL)
	Config.file.Section("meta").Key("description").SetValue(Config.Meta.Description)
	Config.file.Section("meta").Key("title").SetValue(Config.Meta.Title)
	Config.file.Section("meta").Key("language").SetValue(Config.Meta.Language)

	Config.file.Section("database").Key("host").SetValue(Config.Database.Host)
	Config.file.Section("database").Key("port").SetValue(strconv.Itoa(Config.Database.Port))
	Config.file.Section("database").Key("name").SetValue(Config.Database.Name)
	Config.file.Section("database").Key("user").SetValue(Config.Database.User)
	Config.file.Section("database").Key("passwd").SetValue(Config.Database.Passwd)
	Config.file.Section("database").Key("ssl_mode").SetValue(Config.Database.SSLMode)

	Config.file.Section("redis").Key("host").SetValue(Config.Redis.Host)
	Config.file.Section("redis").Key("port").SetValue(strconv.Itoa(Config.Redis.Port))
	Config.file.Section("redis").Key("db").SetValue(strconv.Itoa(Config.Redis.DB))
	Config.file.Section("redis").Key("passwd").SetValue(Config.Redis.Passwd)

	Config.file.Section("http").Key("cookie").SetValue(Config.HTTP.Cookie)
	Config.file.Section("http").Key("api_key").SetValue(Config.HTTP.ApiKey)

	Config.file.Section("cloudflare").Key("email").SetValue(Config.Cloudflare.Email)
	Config.file.Section("cloudflare").Key("api_key").SetValue(Config.Cloudflare.ApiKey)
	Config.file.Section("cloudflare").Key("zone_tag").SetValue(Config.Cloudflare.ZoneTag)

	Config.file.Section("directories").Key("data").SetValue(Config.Directories.Data)

	return Config.file.SaveTo(opts.Path)
}
