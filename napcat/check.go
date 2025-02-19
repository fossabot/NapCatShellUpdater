package napcat

import (
	"fmt"
	"github.com/Sn0wo2/NapCatShellUpdater/flags"
	"github.com/Sn0wo2/NapCatShellUpdater/helper"
	"github.com/Sn0wo2/NapCatShellUpdater/log"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func CheckNapCatUpdate() {
	newVersion := fetchLastNapCatDownloadURL()
	currentVersion := getCurrentNapCatVersion()
	if newVersion != currentVersion {
		log.Info("NapCatShellUpdater", fmt.Sprintf("Updating NapCat from %s to %s", currentVersion, newVersion))
		processAndUpdate(downloadFile(fmt.Sprintf("https://github.com/NapNeko/NapCatQQ/releases/download/%s/NapCat.Shell.zip", newVersion)))
	} else {
		log.Info("NapCatShellUpdater", "NapCat is up to date: ", currentVersion)
	}
}

func ProcessVersionUpdate(ver string) {
	currentVersion := getCurrentNapCatVersion()
	if ver == "" || currentVersion == "" {
		log.Error("NapCatShellUpdater", "Failed to fetch version info", ver, currentVersion)
		return
	}
	if ver != currentVersion {
		processAndUpdate(downloadFile(fmt.Sprintf("https://github.com/NapNeko/NapCatQQ/releases/download/%s/NapCat.Shell.zip", ver)))
	} else {
		log.Info("NapCatShellUpdater", "NapCat is up to date: ", currentVersion)
	}
}

func getCurrentNapCatVersion() (ver string) {
	packageFile, err := os.ReadFile(filepath.Join(flags.Config.Path, "package.json"))
	if err != nil {
		log.Error("NapCatShellUpdater", "failed to read package.json:", err)
		return "v0.0.0(Error)"
	}
	version := gjson.GetBytes(packageFile, "version").String()
	if version == "" {
		version = "0.0.0(Not Found)"
	}
	return "v" + version
}

func fetchLastNapCatDownloadURL() (ver string) {
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/repos/NapNeko/NapCatQQ/releases?per_page=1", nil)
	if err != nil {
		panic(err)
	}

	client := http.DefaultClient
	if flags.Config.Proxy != "" {
		p, err := url.Parse(flags.Config.Proxy)
		if err != nil {
			panic(err)
		}
		client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(p)}}
	}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Error("NapCatShellUpdater", fmt.Sprintf("Failed to fetch version info: %v, status: %d", err, resp.StatusCode))
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	release := gjson.Parse(helper.BytesToString(body)).Array()[0]
	version := release.Get("tag_name").Str
	if version == "" {
		log.Error("NapCatShellUpdater", "Failed to fetch version info\n", helper.BytesToString(body))
		return version
	}
	return version
}
