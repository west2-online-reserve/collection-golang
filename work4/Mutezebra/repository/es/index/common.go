package index

import (
	"context"
	"four/pkg/log"
	"four/repository/es"
	"four/repository/es/model"
)

func IndexCreate(model model.ESModels) (err error) {
	_, err = es.ESClient.CreateIndex(model.Index()).
		BodyString(model.Mapping()).
		Do(context.TODO())
	return err
}

func IndexDelete(model model.ESModels) (err error) {
	_, err = es.ESClient.DeleteIndex(model.Index()).Do(context.TODO())
	return err
}

func IndexExit(model model.ESModels) (exist bool) {
	exist, _ = es.ESClient.IndexExists(model.Index()).Do(context.TODO())
	return
}

func InitIndex() {
	u := &model.User{}
	v := &model.Video{}
	exist := IndexExit(u)
	if !exist {
		err := IndexCreate(u)
		if err != nil {
			log.LogrusObj.Panic(err)
		}
	}
	exist = IndexExit(v)
	if !exist {
		err := IndexCreate(v)
		if err != nil {
			log.LogrusObj.Panic(err)
		}
	}
}
