package main

import (
	"runtime"

	"github.com/Sn0wo2/NapCatShellUpdater/internal/flags"
	"github.com/Sn0wo2/NapCatShellUpdater/pkg/log"
	"github.com/Sn0wo2/NapCatShellUpdater/pkg/napcat"
	"github.com/sirupsen/logrus"
)

func init() {
	flags.InitFlag()

	err := log.InitLogger("", log.DefaultFormatter(), logrus.TraceLevel)
	if err != nil {
		panic(err)
	}

	log.Info("NapCatShellUpdater", `
$$\   $$\                   $$$$$$\           $$\     $$$$$$\ $$\               $$\$$\$$\   $$\               $$\          $$\                       
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

	log.Info("NapCatShellUpdater", "Loading...")

	if runtime.GOOS != "windows" {
		log.Error("NapCatShellUpdater", "Unsupported system:", runtime.GOOS)
	}
}

func main() {
	cv := flags.Config.Version
	if cv == "" {
		napcat.CheckNapCatUpdate()
	} else {
		napcat.ProcessVersionUpdate(cv)
	}
}
