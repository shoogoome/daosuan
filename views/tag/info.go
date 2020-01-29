package tag

import (
	authbase "daosuan/core/auth"
	tagException "daosuan/exceptions/tag"
	"daosuan/models/db"
	paramsUtils "daosuan/utils/params"
	"fmt"
	"github.com/kataras/iris"
)

// 创建tag
func CreateTag(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization) {
	auth.CheckLogin()
	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))

	tag := db.Tag{
		Name: params.Str("name", "名称"),
	}
	db.Driver.FirstOrCreate(&tag, tag)

	ctx.JSON(iris.Map{
		"id": tag.Id,
	})
}

// 获取tag列表，不做分页
func GetTagList(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization) {
	type field struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	var tags []field

	table := db.Driver.Table("tag as t")
	if name := ctx.URLParam("name"); len(name) > 0 {
		table = table.Where("t.name like ?", fmt.Sprintf("%%%s%%", name))
	}
	if course, err := ctx.URLParamInt("course"); err == nil {
		table = table.Table("tag as t, course_tags as ct").Where("ct.course_id = ? and t.id = ct.tag_id", course)
	}
	table.Select("t.id, t.name").Find(&tags)
	ctx.JSON(iris.Map{
		"tags": tags,
	})
}

// 删除tag
func DeleteTag(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, tid int) {
	auth.CheckAdmin()

	var tag db.Tag
	if err := db.Driver.First(&tag, tid).Error; err != nil || tag.Id == 0 {
		panic(tagException.TagIsNotExsits())
	}

	if err := db.Driver.Exec("delete from product_tags where tag_id = ?", tid).Error; err != nil {
		panic(tagException.TagIsNotExsits())
	}
	db.Driver.Delete(&tag)
	ctx.JSON(iris.Map{
		"id": tid,
	})
}
