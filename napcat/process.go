package napcat

import (
	"github.com/Sn0wo2/NapCatShellUpdater/flags"
	"github.com/Sn0wo2/NapCatShellUpdater/log"
	"github.com/shirou/gopsutil/process"

	"os"
	"path/filepath"
	"time"
)

func processAndUpdate(filename string) {
	log.Info("NapCatShellUpdater", "Waiting NapCatWinBootMain.exe process to end...")
	err := WaitForAllProcessesEnd(filepath.Join(flags.Config.Path, "NapCatWinBootMain.exe"))
	if err != nil {
		panic(err)
	}

	log.Info("NapCatShellUpdater", "Clean target directory...")
	err = cleanDirectory(flags.Config.Path, []string{"config", "logs", "quickLoginExample.bat", "update.bat", filepath.Base(os.Args[0]), filename})
	if err != nil {
		log.RPanic(err)
	}

	log.Info("NapCatShellUpdater", "Extracting new version...")
	err = unzipWithExclusion(filename, flags.Config.Path, []string{"quickLoginExample.bat"})
	if err != nil {
		panic(err)
	}

	err = os.Remove(filename)
	if err != nil {
		panic(err)
	}
}

func WaitForProcess(targetPath string) (*process.Process, error) {
	initialProcs, err := process.Processes()
	if err != nil {
		return nil, err
	}
	initialProcIDs := make(map[int32]struct{})
	for _, p := range initialProcs {
		initialProcIDs[p.Pid] = struct{}{}
	}

	for {
		// 减少轮询频率，例如从1秒调整到2秒
		time.Sleep(1 * time.Second)

		procs, err := process.Processes()
		if err != nil {
			return nil, err
		}

		for _, proc := range procs {
			if _, existed := initialProcIDs[proc.Pid]; existed {
				continue
			}

			path, err := proc.Exe()
			if err != nil {
				continue
			}

			if path == targetPath {
				return proc, nil
			}
		}
	}
}

func WaitForAllProcessesEnd(targetPath string) error {
	for {
		allProcs, err := process.Processes()
		if err != nil {
			return err
		}

		activeProcs := 0

		for _, proc := range allProcs {
			exePath, err := proc.Exe()
			if err != nil {
				continue
			}

			if filepath.Dir(exePath) == targetPath {
				activeProcs++
			}
		}

		if activeProcs == 0 {
			return nil
		}

		time.Sleep(2 * time.Second)
	}
}
