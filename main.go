package main

import (
	"fmt"
	"github.com/Sn0wo2/NapCatShellUpdater/flags"
	"github.com/Sn0wo2/NapCatShellUpdater/log"
	"github.com/Sn0wo2/NapCatShellUpdater/napcat"
	"github.com/Sn0wo2/NapCatShellUpdater/napcat/login"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"runtime"
	"time"
)

func init() {
	flags.InitFlag()

	err := log.InitLogger("", log.DefaultFormatter(), logrus.TraceLevel)
	if err != nil {
		panic(err)
	}

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

	if runtime.GOOS != "windows" {
		log.Error("NapCatShellUpdater", "Unsupported system:", runtime.GOOS)
	}
}

func main() {
	napcat.CheckNapCatUpdate()
	if flags.Config.Login {
		log.Info("NapCatShellUpdater", "Wating NapCat process to login...")
		ncProc, err := napcat.WaitForProcess(filepath.Join(flags.Config.Path, "NapCatWinBootMain.exe"))
		if err != nil {
			panic(err)
		}
		log.Debug("NapCatShellUpdater", "NapCat process found:", ncProc.String())
		log.Debug("NapCatShellUpdater", "Waiting 12s to full load NapCat")
		time.Sleep(12 * time.Second)
		log.Info("NapCatShellUpdater", "Login to NapCat Panel...")
		login.LoginNapCat()
	}
}
