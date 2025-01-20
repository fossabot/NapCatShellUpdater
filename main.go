package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"github.com/Sn0wo2/NapCatShellUpdater/helper"
	"github.com/shirou/gopsutil/process"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"time"
)

var (
	path           string
	run            bool
	proxy          string
	debug          bool
	napCatPanelURL string
	napCatToken    string
)

func init() {
	flag.StringVar(&path, "path", "./", "NapCat path")
	flag.BoolVar(&run, "run", true, "Run NapCat")
	flag.StringVar(&proxy, "proxy", "", "HTTP Proxy")
	flag.BoolVar(&debug, "debug", false, "Enable debug logging")
	flag.StringVar(&napCatPanelURL, "napCatPanelURL", "", "NapCat Panel URL")
	flag.StringVar(&napCatToken, "napCatToken", "", "NapCat Token")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println(`$$\   $$\                   $$$$$$\           $$\     $$$$$$\ $$\               $$\$$\$$\   $$\               $$\          $$\                       
$$$\  $$ |                 $$  __$$\          $$ |   $$  __$$\$$ |              $$ $$ $$ |  $$ |              $$ |         $$ |                      
$$$$\ $$ |$$$$$$\  $$$$$$\ $$ /  \__|$$$$$$\$$$$$$\  $$ /  \__$$$$$$$\  $$$$$$\ $$ $$ $$ |  $$ |$$$$$$\  $$$$$$$ |$$$$$$\$$$$$$\   $$$$$$\  $$$$$$\  
$$ $$\$$ |\____$$\$$  __$$\$$ |      \____$$\_$$  _| \$$$$$$\ $$  __$$\$$  __$$\$$ $$ $$ |  $$ $$  __$$\$$  __$$ |\____$$\_$$  _| $$  __$$\$$  __$$\ 
$$ \$$$$ |$$$$$$$ $$ /  $$ $$ |      $$$$$$$ |$$ |    \____$$\$$ |  $$ $$$$$$$$ $$ $$ $$ |  $$ $$ /  $$ $$ /  $$ |$$$$$$$ |$$ |   $$$$$$$$ $$ |  \__|
$$ |\$$$ $$  __$$ $$ |  $$ $$ |  $$\$$  __$$ |$$ |$$\$$\   $$ $$ |  $$ $$   ____$$ $$ $$ |  $$ $$ |  $$ $$ |  $$ $$  __$$ |$$ |$$\$$   ____$$ |      
$$ | \$$ \$$$$$$$ $$$$$$$  \$$$$$$  \$$$$$$$ |\$$$$  \$$$$$$  $$ |  $$ \$$$$$$$\$$ $$ \$$$$$$  $$$$$$$  \$$$$$$$ \$$$$$$$ |\$$$$  \$$$$$$$\$$ |      
\__|  \__|\_______$$  ____/ \______/ \_______| \____/ \______/\__|  \__|\_______\__\__|\______/$$  ____/ \_______|\_______| \____/ \_______\__|      
                  $$ |                                                                         $$ |                                                  
                  $$ |                                                                         $$ |                                                  
                  \__|                                                                         \__|                                                  `)

	checkSystem()
}

func main() {
	newVersion, downloadURL := fetchLastNapCatDownloadURL()
	currentVersion := getCurrentNapCatVersion()
	if newVersion != currentVersion {
		logDebug("Updating NapCat from %s to %s", currentVersion, newVersion)
		processAndUpdate(downloadFile(downloadURL))
	} else {
		logDebug("NapCat is up to date: %s", currentVersion)
	}

	if run {
		err := exec.Command("cmd", "/C", "start", "", "quickLoginExample.bat").Run()
		if err != nil {
			log.Fatalf("Failed to run NapCat: %v", err)
		}
		time.Sleep(12 * time.Second)
		token := loginNapCatPanel()
		setNapCatQuickLogin(token, getNapCatPanelLoginList(token))
	}

}
func loginNapCatPanel() (token string) {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/NapCat/api/auth/login", napCatPanelURL), strings.NewReader(fmt.Sprintf(`{"token":"%s"}`, napCatToken)))
	if err != nil {
		log.Fatalf("Failed to create HTTP request: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to login: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}
	return gjson.Parse(helper.BytesToString(body)).Get("data.Credential").String()
}

func getNapCatPanelLoginList(token string) []int64 {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/NapCat/api/QQLogin/GetQuickLoginList", napCatPanelURL), nil)
	if err != nil {
		log.Fatalf("Failed to create HTTP request: %v", err)
	}
	req.Header.Add("Authorization", "Bearer "+token)
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to login: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}
	var loginList []int64
	gjson.Parse(helper.BytesToString(body)).Get("data").ForEach(func(key, value gjson.Result) bool {
		loginList = append(loginList, value.Int())
		return true
	})
	return loginList
}

func setNapCatQuickLogin(token string, loginList []int64) {
	if len(loginList) < 1 {
		return
	}
	for _, uin := range loginList {
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/NapCat/api/QQLogin/SetQuickLogin", napCatPanelURL), strings.NewReader(fmt.Sprintf(`{"uin":"%d"}`, uin)))
		if err != nil {
			log.Fatalf("Failed to create HTTP request: %v", err)
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+token)
		client := http.DefaultClient
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Failed to login: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			logDebug("Quick login: %d", uin)
		} else {
			fmt.Println("Failed login fucking napcat:", uin, "status:", resp.StatusCode)
		}
		time.Sleep(222 * time.Millisecond)
	}
}

func checkSystem() {
	if runtime.GOOS != "windows" {
		log.Fatalf("Unsupported system: %s", runtime.GOOS)
	}
}

func fetchLastNapCatDownloadURL() (string, string) {
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/repos/NapNeko/NapCatQQ/releases?per_page=1", nil)
	if err != nil {
		log.Fatalf("Failed to create HTTP request: %v", err)
	}

	client := http.DefaultClient
	if proxy != "" {
		p, err := url.Parse(proxy)
		if err != nil {
			log.Fatalf("Invalid proxy URL: %v", err)
		}
		client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(p)}}
	}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to fetch version info: %v, status: %d", err, resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	release := gjson.ParseBytes(body).Array()[0]
	version := release.Get("tag_name").String()
	downloadURL := release.Get("assets.#(name==NapCat.Shell.zip).browser_download_url").String()
	if downloadURL == "" {
		log.Fatalf("Download URL not found for version %s", version)
	}
	return version, downloadURL
}

func getCurrentNapCatVersion() string {
	data, err := os.ReadFile(filepath.Join(path, "package.json"))
	if err != nil {
		log.Fatalf("Failed to read package.json: %v", err)
	}
	version := gjson.GetBytes(data, "version").String()
	return "v" + version
}

func downloadFile(downloadURL string) string {
	req, err := http.NewRequest(http.MethodGet, downloadURL, nil)
	if err != nil {
		log.Fatalf("Failed to create download request: %v", err)
	}

	client := http.DefaultClient
	if proxy != "" {
		p, err := url.Parse(proxy)
		if err != nil {
			log.Fatalf("Invalid proxy URL: %v", err)
		}
		client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(p)}}
	}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to download file: %v, status: %d", err, resp.StatusCode)
	}
	defer resp.Body.Close()

	filename := fmt.Sprintf("NapCat.Shell(%s).zip", time.Now().Format("20060102150405"))
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}
	return filename
}

func processAndUpdate(filename string) {
	processes, err := process.Processes()
	if err != nil {
		log.Fatalf("Failed to list processes: %v", err)
	}

	ncb := "NapCatWinBootMain.exe"
	for _, p := range processes {
		name, _ := p.Name()
		if name == ncb {
			exe, _ := p.Exe()
			// 确保进程属于目标目录
			absncb, err := filepath.Abs(path + "/" + ncb)
			if err != nil {
				log.Fatalf("Failed to get absolute path of %s: %v", ncb, err)
			}
			if exe == absncb {
				if err := p.Kill(); err != nil {
					logDebug("Failed to kill process %s: %v", ncb, err)
				} else {
					logDebug("Killed process: %s", exe)
				}
				// 等待dll占用被解除
				time.Sleep(100 * time.Millisecond)
			}
		}
	}

	logDebug("Cleaning target directory...")
	err = cleanDirectory(path, []string{"config", "logs", "quickLoginExample.bat", "update.bat", filepath.Base(os.Args[0]), filename})
	if err != nil {
		log.Fatalf("Failed to clean target directory: %v", err)
	}

	logDebug("Extracting new version...")
	err = unzipWithExclusion(filename, path, []string{"quickLoginExample.bat"})
	if err != nil {
		log.Fatalf("Failed to extract new version: %v", err)
	}

	err = os.Remove(filename)
	if err != nil {
		logDebug("Failed to remove file %s: %v", filename, err)
	}
}

func cleanDirectory(targetPath string, exclude []string) error {
	files, err := os.ReadDir(targetPath)
	if err != nil {
		return fmt.Errorf("failed to read target directory: %w", err)
	}

	for _, file := range files {
		// 跳过当前执行的文件
		if slices.Contains(exclude, file.Name()) {
			continue
		}
		filePath := filepath.Join(targetPath, file.Name())
		if file.IsDir() {
			err := os.RemoveAll(filePath)
			if err != nil {
				return fmt.Errorf("failed to remove directory %s: %w", filePath, err)
			}
		} else {
			err := os.Remove(filePath)
			if err != nil {
				return fmt.Errorf("failed to remove file %s: %w", filePath, err)
			}
		}
	}
	return nil
}

func unzipWithExclusion(src, dest string, exclude []string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("failed to open zip: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		skip := false
		for _, e := range exclude {
			if f.Name == e {
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			os.MkdirAll(filepath.Dir(fpath), os.ModePerm)
			outFile, err := os.Create(fpath)
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}
			rc, err := f.Open()
			if err != nil {
				outFile.Close()
				return fmt.Errorf("failed to open zip entry: %w", err)
			}
			_, err = io.Copy(outFile, rc)
			outFile.Close()
			rc.Close()
			if err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}
		}
	}
	return nil
}

func logDebug(format string, v ...interface{}) {
	if debug {
		log.Printf("[DEBUG] "+format, v...)
	}
}
