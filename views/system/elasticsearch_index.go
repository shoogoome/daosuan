package system

import (
	"bytes"
	"daosuan/core/auth"
	"daosuan/core/elasticsearch"
	"daosuan/core/queue"
	"daosuan/models/db"
	"encoding/json"
	"github.com/kataras/iris"
)

// 重建elasticsearch索引
func ResetElasticsearchIndex(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization) {
	//auth.CheckAdmin()

	queue.Task <- func() {
		var products []db.Product
		db.Driver.Find(&products)

		for _, product := range products {
			if re, err := json.Marshal(product); err == nil {
				elasticsearch.Create("product", "root", product.Id, bytes.NewBuffer(re))
			}
		}
	}
	ctx.JSON(iris.Map {
		"status": "重建任务已加入后台队列。详情可查看日志",
	})
}
