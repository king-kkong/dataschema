package main

import (
	"fmt"

	"github.com/pulingfu/tblschema/cmd/test_dataer/model"
	"github.com/pulingfu/tblschema/dvaplugin"
	"github.com/tidwall/gjson"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	var configs = &gorm.Config{}
	db, err := gorm.Open(mysql.Open("root:tiger@(127.0.0.1:3306)/pulingfu?charset=utf8mb4&parseTime=True&loc=Local"), configs)
	if err != nil {
		panic(err)
	}

	var companys []model.TblCompany
	db.Model(&companys).Where("id in ?", []int{3, 1}).Find(&companys)

	//加载关系
	rl := dvaplugin.NewRelationLoader(companys, true)
	//加载公司创建者  使用默认比较方法和默认的 db查询关系
	rl.AddRelation(dvaplugin.HAS_ONE, "creator", "creator_id", "id", model.TblUser{}, nil, nil)
	//加载项目
	rl.AddRelation(dvaplugin.HAS_MANEY, "projects", "id", "company_id", model.TblProject{}, func(p, s gjson.Result) bool {
		return p.Get("id").Int() == s.Get("company_id").Int() && p.Get("id").Int() > 0
	}, nil)
	//加载项目创建者
	rl.AddRelation(dvaplugin.HAS_ONE, "projects.creator", "creator_id", "id", model.TblUser{}, nil,
		db.Model(&model.TblUser{}).Where("id = ?", 1).Limit(1))
	//加载项目细节
	rl.AddRelation(dvaplugin.HAS_MANEY, "projects.project_details", "id", "project_id", model.TblProjectDetail{}, nil, nil)

	//错误案例
	rl.AddRelation(dvaplugin.HAS_MANEY, "xxx.projects.project_details", "id", "project_id", model.TblProject{}, nil, nil)
	rl.LoadResult(db)

	jv := dvaplugin.VtoJson(rl.GetResult())
	fmt.Println(jv.String())
}
