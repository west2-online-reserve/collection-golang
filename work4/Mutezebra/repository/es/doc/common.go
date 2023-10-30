package doc

import (
	"context"
	"encoding/json"
	"four/repository/es"
	"four/repository/es/model"
	"four/types"
	"github.com/olivere/elastic/v7"
	"math"
	"strconv"
)

func DocCreate(model model.ESModels) error {
	_, err := es.ESClient.Index().Index(model.Index()).BodyJson(&model).Do(context.TODO())
	return err
}

func DocUpdate(model model.ESModels, filed string, value interface{}, id string) error {
	_, err := es.ESClient.Update().Index(model.Index()).Doc(map[string]any{
		filed: value,
	}).Id(id).Refresh("true").Do(context.TODO())
	return err
}

func DocSearch(model model.ESModels, filed string, id uint) (*elastic.SearchHit, string, error) {
	query := elastic.NewTermQuery(filed, id)

	resp, err := es.ESClient.Search(model.Index()).Query(query).Size(1).Do(context.TODO())
	if err != nil {
		return nil, "", err
	}
	if len(resp.Hits.Hits) != 1 {
		return nil, "", RecordNotExist
	}
	result := resp.Hits.Hits[0]
	return result, result.Id, nil
}

// TODO 去重,分页
func Search(u *model.User, req *types.SearchReq) ([]uint, error) {
	v := &model.Video{}
	// 尝试查询精确人名
	vid1 := make([]uint, 0) // 存放找到的视频的id
	if req.Pages == 1 {     // 如果页数为1再找用户
		if u != nil { // 如果u不为空的话就说明查找到了精确人名，然后根据uid去找作者的视频
			query := elastic.NewBoolQuery()
			condition := elastic.NewTermsQuery("uid", u.Uid)
			query.Must(condition)
			resp, err := es.ESClient.Search(v.Index()).Query(query).Do(context.TODO())
			if err != nil {
				return nil, err
			}
			if len(resp.Hits.Hits) == 0 { // 如果这个长度为0说明这个用户没有发过作品
			}

			if len(resp.Hits.Hits) <= 5 {
				for _, value := range resp.Hits.Hits {
					json.Unmarshal(value.Source, v)
					vid1 = append(vid1, v.Vid)
				}
			}

			if len(resp.Hits.Hits) > 5 { // 即便发了很多作品这里也只匹配5条
				for i := 0; i < 5; i++ {
					json.Unmarshal(resp.Hits.Hits[i].Source, v)
					vid1 = append(vid1, v.Vid)
				}
			}
		}
	}

	// 找完精确人名相关的视频，就来匹配其他的信息
	vid2 := make([]uint, 0)
	query := elastic.NewBoolQuery()
	conditionTitle := elastic.NewMatchQuery("title", req.Content)
	conditionIntro := elastic.NewMatchQuery("intro", req.Content)
	query.Should(conditionTitle, conditionIntro)
	resp, err := es.ESClient.Search(v.Index()).Query(query).Do(context.TODO())
	if err != nil {
		return nil, err
	}
	if len(resp.Hits.Hits)+len(vid1) > 10 { // 如果总数据大于10条，那么就只取10条
		for i := 0; i < 10-len(vid1); i++ {
			err = json.Unmarshal(resp.Hits.Hits[i].Source, v)
			if err != nil {
				return nil, err
			}
			vid2 = append(vid2, v.Vid)
		}
	} else { // 如果总数据不超过10条，那么就把这个resp中的所有视频id都读取
		for _, value := range resp.Hits.Hits {
			err = json.Unmarshal(value.Source, v)
			if err != nil {
				return nil, err
			}
			vid2 = append(vid2, v.Vid)
		}
	}

	vid1 = append(vid1, vid2...) // 将数据全都放在vid1中
	return vid1, nil
}

func VideoFilter(start, end int64, conds string, req *types.FilterReq) ([]uint, error) {
	v := &model.Video{}
	Pages, _ := strconv.Atoi(req.Pages)
	Size, _ := strconv.Atoi(req.Size)
	viewS, _ := strconv.ParseInt(req.ViewStart, 10, 64)
	viewE, _ := strconv.ParseInt(req.ViewEnd, 10, 64)
	if viewS < 0 || viewS > math.MaxInt32 {
		viewS = 0
	}
	if viewE < 0 || viewE > math.MaxInt32 {
		viewE = math.MaxInt32
	}

	from := (Pages - 1) * Size

	query := elastic.NewBoolQuery()

	timeFilter := elastic.NewRangeQuery("create_timestamp").Gte(start).Lte(end)
	viewFilter := elastic.NewRangeQuery("view").Gte(viewS).Lte(viewE)

	if conds != "" {
		tagFilter := elastic.NewMatchQuery("tag", conds)
		query.Must(tagFilter)
	}

	query.Must(timeFilter, viewFilter)

	tagIdsQuery := elastic.NewMatchQuery("tag", conds)
	titleMatchQuery := elastic.NewMatchQuery("title", req.Content)
	introMatchQuery := elastic.NewMatchQuery("intro", req.Content)

	query.Should(tagIdsQuery, titleMatchQuery, introMatchQuery)

	resp, err := es.ESClient.Search(v.Index()).Query(query).From(from).Size(Size).Do(context.TODO())
	vid := make([]uint, 0)
	if err != nil {
		return nil, err
	}
	if len(resp.Hits.Hits) == 0 {
		return vid, nil
	}
	for _, value := range resp.Hits.Hits {
		json.Unmarshal(value.Source, v)
		vid = append(vid, v.Vid)
	}
	return vid, nil
}
