package main

import (
	"fmt"
	"github.com/Sn0wo2/NapCatShellUpdater/flags"
	"github.com/Sn0wo2/NapCatShellUpdater/log"
	"github.com/Sn0wo2/NapCatShellUpdater/napcat"
	"github.com/sirupsen/logrus"
	"runtime"
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
