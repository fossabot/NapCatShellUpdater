package flags

import (
	"flag"
	"path/filepath"

	"github.com/Sn0wo2/NapCatShellUpdater/pkg/log"
)

var Config struct {
	Path        string
	Version     string
	Proxy       string
	DownloadURL string
	Exclude     string
	Debug       bool
}

func InitFlag() bool {
	path := "./"
	flag.StringVar(&path, "path", "./", "NapCat path")
	flag.StringVar(&Config.Version, "version", "", "Update NapCat Version")
	flag.StringVar(&Config.Proxy, "proxy", "", "HTTP Proxy")
	flag.StringVar(&Config.DownloadURL, "download-url", "", "Download NapCat URL")
	flag.StringVar(&Config.Exclude, "exclude", "", "Exclude files")
	flag.BoolVar(&Config.Debug, "debug", true, "Enable debug logging")
	flag.Parse()
	var err error
	Config.Path, err = filepath.Abs(path)
	if err != nil {
		Config.Path = path
		log.RPanic(err)
	}
	return true
}
