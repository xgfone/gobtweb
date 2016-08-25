package g

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/xgfone/gobtweb/logger"
	"gopkg.in/olivere/elastic.v3"
)

var (
	ElasticClient *elastic.Client
	// Repository    repository.Repository
	Conf   Config
	Logger *logger.Logger
)

type Config struct {
	Address  string `json:"address"`
	Elastic  string `json:"elastic"`
	Database string `json:"db"`
	LogFile  string `json:"logfile"`
	LogLevel string `json:"loglevel"`
	PageSize int    `json:"page_size"`

	TemplateDirectory string `json:"template_directory"`
	StaticDirectory   string `json:"static_directory"`
}

func initConfig(filename string) {
	if f, err := os.Open(filename); err != nil {
		panic(err)
	} else if data, err := ioutil.ReadAll(f); err != nil {
		panic(err)
	} else if err = json.Unmarshal(data, &Conf); err != nil {
		panic(err)
	}
}

func Init(config_file string) {
	var err error
	initConfig(config_file)

	// if Repository, err = repository.NewMysqlRepository(Conf.Database, 256, 256); err != nil {
	// 	panic(err)
	// }

	if ElasticClient, err = elastic.NewClient(elastic.SetURL(Conf.Elastic)); err != nil {
		panic(err)
	} else if exists, err := ElasticClient.IndexExists("torrent").Do(); err != nil {
		panic(err)
	} else if !exists {
		if _, err := ElasticClient.CreateIndex("torrent").Do(); err != nil {
			panic(err)
		}
	}

	if Logger, err = logger.NewLogger(Conf.LogLevel, Conf.LogFile); err != nil {
		panic(err)
	}

	if Conf.PageSize < 1 {
		Conf.PageSize = 20
	}
}
