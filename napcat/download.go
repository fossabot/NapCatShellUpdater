package napcat

import (
	"archive/zip"
	"fmt"
	"github.com/Sn0wo2/NapCatShellUpdater/flags"
	"github.com/Sn0wo2/NapCatShellUpdater/log"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"time"
)

func downloadFile(downloadURL string) string {
	req, err := http.NewRequest(http.MethodGet, downloadURL, nil)
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

	filename := fmt.Sprintf("NapCat.Shell(%d).zip", time.Now().Unix())
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		panic(err)
	}
	return filename
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
				log.Info("NapCatShellUpdater", "Waiting QQ.exe process to end...")
				err = <-WaitForAllProcessesEnd("QQ.exe", false)
				if err != nil {
					panic(err)
				}
				return fmt.Errorf("failed to remove directory %s: %w", filePath, err)
			}
		} else {
			err := os.Remove(filePath)
			if err != nil {
				log.Info("NapCatShellUpdater", "Waiting QQ.exe process to end...")
				err2 := <-WaitForAllProcessesEnd("QQ.exe", false)
				if err2 != nil {
					panic(err2)
				}
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
