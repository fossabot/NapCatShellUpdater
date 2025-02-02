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
	newVersion, downloadURL := fetchLastNapCatDownloadURL()
	currentVersion := getCurrentNapCatVersion()
	if newVersion != currentVersion {
		log.Info("NapCatShellUpdater", fmt.Sprintf("Updating NapCat from %s to %s", currentVersion, newVersion))
		processAndUpdate(downloadFile(downloadURL))
	} else {
		log.Info("NapCatShellUpdater", "NapCat is up to date: ", currentVersion)
	}
}

func getCurrentNapCatVersion() string {
	data, err := os.ReadFile(filepath.Join(flags.Config.Path, "package.json"))
	if err != nil {
		panic(err)
	}
	version := gjson.GetBytes(data, "version").String()
	if version == "" {
		version = "0.0.0(Not Found)"
	}
	return "v" + version
}

func fetchLastNapCatDownloadURL() (string, string) {
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
	version := release.Get("tag_name").String()
	downloadURL := release.Get("assets.#(name==NapCat.Shell.zip).browser_download_url").String()
	if version == "" || downloadURL == "" {
		log.Error("NapCatShellUpdater", "Failed to fetch version info\n", helper.BytesToString(body))
		os.Exit(1)
	}
	return version, downloadURL
}
