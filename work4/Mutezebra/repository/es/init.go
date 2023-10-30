package es

import (
	"crypto/tls"
	"fmt"
	"four/config"
	"four/pkg/log"
	"github.com/olivere/elastic/v7"
	log2 "log"
	"net/http"
)

var ESClient *elastic.Client

func InitES() {
	conf := config.Config.ES
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // 忽略证书验证
		},
	}
	httpClient := &http.Client{Transport: tr}
	es, err := elastic.NewClient(
		elastic.SetURL(conf.Address),
		elastic.SetBasicAuth(conf.UserName, conf.Password),
		elastic.SetSniff(false), // 关闭嗅探器,不然会被检验不通过
		elastic.SetHttpClient(httpClient),
	)
	if err != nil {
		log2.Println(err)
		fmt.Println(err)
		log.LogrusObj.Panic("ESClient connect failed,error: ", err)
	}
	log.LogrusObj.Infoln("es client connect success")
	ESClient = es
}
