package flags

import (
	"flag"
	"github.com/Sn0wo2/NapCatShellUpdater/log"
	"path/filepath"
)

var Config struct {
	Path    string
	Version string
	Proxy   string
	Debug   bool
}

func InitFlag() bool {
	path := "./"
	flag.StringVar(&path, "path", "./", "NapCat path")
	var err error
	Config.Path, err = filepath.Abs(path)
	if err != nil {
		Config.Path = path
		log.RPanic(err)
	}
	flag.StringVar(&Config.Version, "version", "", "Update NapCat Version")
	flag.StringVar(&Config.Proxy, "proxy", "", "HTTP Proxy")
	flag.BoolVar(&Config.Debug, "debug", true, "Enable debug logging")
	flag.Parse()
	return true
}
