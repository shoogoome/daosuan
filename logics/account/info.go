package accountLogic

import (
	"daosuan/constants"
	authbase "daosuan/core/auth"
	"daosuan/core/cache"
	"daosuan/exceptions/account"
	"daosuan/logics/resource"
	"daosuan/models/db"
	"daosuan/models/dto"
	"daosuan/utils/hash"
	paramsUtils "daosuan/utils/params"
	"encoding/json"
)

var field = []string{
	"Nickname", "Email", "Id", "Role", "Phone", "PhoneValidated", "UpdateTime",
	"EmailValidated", "Avator", "Motto", "CreateTime", "Init",
}

type follow struct {
	Id       int    `json:"id"`
	Nickname string `json:"nickname"`
	Motto    string `json:"motto"`
}

type AccountLogic struct {
	auth    authbase.DaoSuanAuthAuthorization
	account db.Account
}

func NewAccountLogic(auth authbase.DaoSuanAuthAuthorization, aid ...int) AccountLogic {
	var account db.Account

	if len(aid) > 0 {
		if err := db.Driver.GetOne("account", aid[0], &account); err != nil || account.Id == 0 {
			panic(accountException.AccountIsNotExists())
		}
	} else {
		account = *auth.AccountModel()
	}
	return AccountLogic{
		account: account,
		auth:    auth,
	}
}

func (a *AccountLogic) SetAccountModel(account db.Account) {
	a.account = account
}

func (a *AccountLogic) AccountModel() *db.Account {
	return &a.account
}

func (a *AccountLogic) GetAccountInfo() interface{} {

	if len(a.account.Avator) > 0 {
		a.account.Avator = resourceLogic.GenerateToken(a.account.Avator, -1, -1)
	}

	return paramsUtils.ModelToDict(a.account, field)
}

// 获取关注的人
func (a *AccountLogic) GetFollowing() []follow {
	var follow []follow
	// 读取缓存
	if payload, err := cache.Dijan.Get(paramsUtils.CacheBuildKey(constants.FollowingModel, a.account.Id)); err == nil {
		if err = json.Unmarshal(payload, &follow); err == nil {
			return follow
		}
	}

	db.Driver.
		Table("account_follow as f, account as a").
		Where("f.source_id = ? and a.id = f.target_id", a.account.Id).
		Select("a.id, a.nickname, a.motto").
		Find(&follow)
	if payload, err := json.Marshal(follow); err == nil {
		v := hash.RandInt64(240, 240*5)
		cache.Dijan.Set(paramsUtils.CacheBuildKey(constants.FollowingModel, a.account.Id), payload, int(v)*60*60)
	}
	return follow
}

// 获取关注的我的人
func (a *AccountLogic) GetFollowers() []follow {
	var follow []follow
	// 读取缓存
	if payload, err := cache.Dijan.Get(paramsUtils.CacheBuildKey(constants.FollowerModel, a.account.Id)); err == nil {
		if err = json.Unmarshal(payload, &follow); err == nil {
			return follow
		}
	}

	db.Driver.
		Table("account_follow as f, account as a").
		Where("f.target_id = ? and a.id = f.source_id", a.account.Id).
		Select("a.id, a.nickname, a.motto").
		Find(&follow)
	if payload, err := json.Marshal(follow); err == nil {
		v := hash.RandInt64(240, 240*5)
		cache.Dijan.Set(paramsUtils.CacheBuildKey(constants.FollowerModel, a.account.Id), payload, int(v)*60*60)
	}
	return follow
}

// 获取star产品
func (a *AccountLogic) GetStars() []dto.ProductList {

	var stars []dto.ProductList
	db.Driver.
		Table("account_star as s, product as p").
		Where("s.account_id = ? and p.id = s.product_id", a.account.Id).
		Select("p.id, p.update_time, p.cover, p.create_time, p.description, p.name, p.status, p.star").
		Order("s.create_time desc").Find(&stars)
	// star产品状态需实时，不适合缓存
	return stars
}

// 获取步伐的产品
func (a *AccountLogic) GetProduct() []dto.ProductList {
	var lists []dto.ProductList
	// 读取缓存
	if payload, err := cache.Dijan.Get(paramsUtils.CacheBuildKey(constants.AccountProductModel, a.account.Id)); err == nil && payload != nil {
		if err = json.Unmarshal(payload, &lists); err == nil {
			return lists
		}
	}

	db.Driver.Table("product as p").
		Where("author_id = ?", a.account.Id).
		Select("p.id, p.update_time, p.cover, p.create_time, p.description, p.name, p.status, p.star").
		Order("p.update_time desc").Find(&lists)
	for i := 0; i < len(lists); i++ {
		lists[i].Cover = resourceLogic.GenerateToken(lists[i].Cover, -1, -1)
	}

	if payload, err := json.Marshal(lists); err == nil {
		v := hash.RandInt64(240, 240*5)
		cache.Dijan.Set(paramsUtils.CacheBuildKey(constants.AccountProductModel, a.account.Id), payload, int(v) * 60 * 60)
	}
	return lists
}


// 检测昵称是否存在
func IsNicknameExists(nickname string) bool {
	if name, err := cache.Dijan.Get(paramsUtils.CacheBuildKey(constants.NicknameModel, nickname)); err == nil && name != nil {
		return true
	}

	var account db.Account
	if err := db.Driver.Where("nickname = ?", nickname).First(&account).Error; err != nil || account.Id == 0 {
		return false
	}
	v := hash.RandInt64(240, 240*5)
	cache.Dijan.Set(paramsUtils.CacheBuildKey(constants.NicknameModel, nickname), []byte(nickname), int(v)*60*60)
	return true
}
