package napcat

import (
	"time"

	"github.com/shirou/gopsutil/process"
)

func WaitForProcess(targetPath string) (<-chan *process.Process, <-chan error) {
	resultChan := make(chan *process.Process)
	errorChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		defer close(errorChan)

		initialProcs, err := process.Processes()
		if err != nil {
			errorChan <- err
			return
		}

		initialProcIDs := make(map[int32]struct{})
		for _, p := range initialProcs {
			initialProcIDs[p.Pid] = struct{}{}
		}
		for {
			time.Sleep(2 * time.Second)

			procs, err := process.Processes()
			if err != nil {
				errorChan <- err
				return
			}

			for _, proc := range procs {
				if _, existed := initialProcIDs[proc.Pid]; !existed {
					path, err := proc.Exe()
					if err != nil {
						continue
					}
					if path == targetPath {
						resultChan <- proc
						return
					}
				}
			}
		}
	}()

	return resultChan, errorChan
}

func WaitForAllProcessesEnd(target string, abs bool) <-chan error {
	errorChan := make(chan error, 1)

	go func() {
		defer close(errorChan)

		for {
			allProcs, err := process.Processes()
			if err != nil {
				errorChan <- err
				return
			}

			activeProcs := 0
			for _, proc := range allProcs {
				exePath, err := proc.Exe()
				if err != nil {
					continue
				}

				if abs {
					if exePath == target {
						activeProcs++
					}
				} else {
					if name, err := proc.Name(); err != nil && name == target {
						activeProcs++
					}
				}
			}

			if activeProcs == 0 {
				return
			}

			time.Sleep(2 * time.Second)
		}
	}()

	return errorChan
}
