package napcat

import (
	"time"

	"github.com/shirou/gopsutil/process"
)

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
