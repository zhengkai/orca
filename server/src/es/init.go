package es

import (
	"project/config"
	"project/zj"

	"github.com/elastic/go-elasticsearch/v8"
)

var theClient *elasticsearch.Client

// Init ...
func Init() (err error) {

	theClient, err = elasticsearch.NewClient(elasticsearch.Config{
		Username: config.ESUser,
		Password: config.ESPass,
	})
	if err != nil {
		return err
	}

	res, err := theClient.Info()
	if err != nil {
		return err
	}
	defer res.Body.Close()

	zj.J(`elasticsearch`, res.String())

	return
}
