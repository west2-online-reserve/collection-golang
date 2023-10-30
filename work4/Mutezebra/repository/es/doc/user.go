package doc

import (
	"context"
	"errors"
	"four/repository/es"
	"four/repository/es/model"
	"github.com/olivere/elastic/v7"
)

var (
	RecordNotExist = errors.New("record not exist")
)

func SearchUser(userName string) (hit *elastic.SearchHit, err error) {
	user := &model.User{}
	query := elastic.NewBoolQuery()
	must := elastic.NewTermQuery("user_name.keyword", userName)
	query.Must(must)
	resp, err := es.ESClient.Search(user.Index()).Query(query).Size(1).Do(context.TODO())
	if err != nil {
		return nil, err
	}
	if len(resp.Hits.Hits) == 0 {
		return nil, RecordNotExist
	}
	return resp.Hits.Hits[0], nil
}
