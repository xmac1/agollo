package agollo

import (
	"io"
	"log"
	"os"
)

//start apollo
func Start() error {
	return StartWithLogger(nil)
}

var logger *log.Logger

func StartWithLogger(writer io.Writer) error {
	if writer == nil {
		writer = os.Stdout
	}

	logger.Println(writer, "", log.LstdFlags | log.Lshortfile)

  //init server ip list
  go initServerIpList()

	//first sync
	err := notifySyncConfigServices()

	//first sync fail then load config file
	if err !=nil{
		config, _ := loadConfigFile(appConfig.BackupConfigPath)
		if config!=nil{
			updateApolloConfig(config,false)
		}
	}

	//start long poll sync config
	go StartRefreshConfig(&NotifyConfigComponent{})

	logger.Println("agollo start finished , error:",err)
	
	return err
}
