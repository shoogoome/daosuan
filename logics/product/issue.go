package productLogic

import (
	authbase "daosuan/core/auth"
	"daosuan/enums/product"
	productException "daosuan/exceptions/product"
	"daosuan/models/db"
)

type IssueLogic struct {
	auth    authbase.DaoSuanAuthAuthorization
	product db.Product
	issue   db.Issue
}

func NewIssueLogic(auth authbase.DaoSuanAuthAuthorization, pid int, iid ...int) IssueLogic {
	var issue db.Issue
	var product db.Product

	if err := db.Driver.GetOne("product", pid, &product); err != nil {
		panic(productException.ProductIsNotExists())
	}
	if product.Status != productEnums.StatusReleased {
		panic(productException.NoPermission())
	}

	if len(iid) > 0 {
		if err := db.Driver.GetOne("issue", iid[0], &issue); err != nil || issue.Id == 0 || issue.ProductId != pid {
			panic(productException.IssueIsNotExists())
		}
	}
	return IssueLogic{
		auth:    auth,
		issue:   issue,
		product: product,
	}
}

// 设置model实体
func (i *IssueLogic) SetIssueModel(issue db.Issue) {
	i.issue = issue
}

// 获取model实体
func (i *IssueLogic) IssueModel() *db.Issue {
	return &i.issue
}
