package elasticsearch

import (
	"log"
	"os"

	"github.com/olivere/elastic/v7"
)

var (
	elasticClient *elastic.Client
)

func Init() {
	var err error
	elasticClient, err = elastic.NewClient(
		elastic.SetURL(`http://es01:9200`),
		elastic.SetSniff(false),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)),
		elastic.SetErrorLog(log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)),
	)

	newVideoIndex()

	if err != nil {
		panic(err)
	}

}
