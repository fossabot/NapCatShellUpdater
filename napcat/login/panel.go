package login

import (
	"bufio"
	"fmt"
	"github.com/Sn0wo2/NapCatShellUpdater/flags"
	"github.com/Sn0wo2/NapCatShellUpdater/helper"
	"github.com/Sn0wo2/NapCatShellUpdater/log"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

func NapCatLogin() {
	token := loginNapCatPanel()
	if token != "" {
		loginList := getNapCatPanelLoginList(token)
		if len(loginList) <= 0 {
			return
		}
		setNapCatQuickLogin(token, loginList)
	}
}

func GetNapCatPanelURLInLogs(dirPath string) (string, string, error) {
	// 验证目录路径
	fileInfo, err := os.Stat(dirPath)
	if err != nil || !fileInfo.IsDir() {
		return "", "", fmt.Errorf("invalid directory path: %s", dirPath)
	}

	// 读取目录条目
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to read directory: %v", err)
	}

	// 预编译正则表达式
	urlTokenRegex := regexp.MustCompile(`(https?://[^\s:/]+:\d+)/webui\?token=([^\s]+)`)

	// 收集并排序日志文件
	var logFiles []struct {
		Path    string
		ModTime time.Time
	}
	for _, entry := range entries {
		if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") || strings.ToLower(filepath.Ext(entry.Name())) != ".log" {
			continue
		}

		fullPath := filepath.Join(dirPath, entry.Name())
		fileInfo, err := entry.Info()
		if err != nil {
			continue
		}

		logFiles = append(logFiles, struct {
			Path    string
			ModTime time.Time
		}{
			Path:    fullPath,
			ModTime: fileInfo.ModTime(),
		})
	}

	if len(logFiles) == 0 {
		return "", "", fmt.Errorf("no log files found in %s", dirPath)
	}

	// 按修改时间排序（最新优先）
	sort.Slice(logFiles, func(i, j int) bool {
		return logFiles[i].ModTime.After(logFiles[j].ModTime)
	})

	// 检查每个日志文件
	for _, logFile := range logFiles {
		f, err := os.Open(logFile.Path)
		if err != nil {
			continue
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			matches := urlTokenRegex.FindStringSubmatch(scanner.Text())
			if len(matches) >= 3 {
				return matches[1], matches[2], nil
			}
		}
	}

	return "", "", fmt.Errorf("no matching URL found in %s", dirPath)
}

func loginNapCatPanel() (token string) {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/auth/login", flags.Config.NapCatPanelURL), strings.NewReader(fmt.Sprintf(`{"token":"%s"}`, flags.Config.NapCatToken)))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	log.Trace("NapCatShellUpdater", "Login to NapCat Panel:", helper.BytesToString(body))
	return gjson.Parse(helper.BytesToString(body)).Get("data.Credential").String()
}

func getNapCatPanelLoginList(token string) []int64 {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/QQLogin/GetQuickLoginList", flags.Config.NapCatPanelURL), nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+token)
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var loginList []int64
	gjson.Parse(helper.BytesToString(body)).Get("data").ForEach(func(key, value gjson.Result) bool {
		loginList = append(loginList, value.Int())
		return true
	})
	log.Trace("NapCatShellUpdater", "Get NapCat Panel Login List:", loginList)
	return loginList
}

func setNapCatQuickLogin(token string, loginList []int64) {
	for _, uin := range loginList {
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/QQLogin/SetQuickLogin", flags.Config.NapCatPanelURL), strings.NewReader(fmt.Sprintf(`{"uin":"%d"}`, uin)))
		if err != nil {
			panic(err)
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+token)
		client := http.DefaultClient
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		log.Trace("NapCatShellUpdater", uin, " | Set NapCat Panel Quick Login:", helper.BytesToString(body))
		time.Sleep(222 * time.Millisecond)
	}
}
