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
	"EmailValidated", "Avator", "Motto", "CreateTime",
}

type AccountLogic interface {
	GetAccountInfo() interface{}
	AccountModel() *db.Account
	SetAccountModel(account db.Account)
	GetFollowers() []follow
	GetFollowing() []follow
	GetStars() []dto.ProductList
}

type follow struct {
	Id       int    `json:"id"`
	Nickname string `json:"nickname"`
	Motto    string `json:"motto"`
}

type accountStruct struct {
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
	return &accountStruct{
		account: account,
		auth:    auth,
	}
}

func (a *accountStruct) SetAccountModel(account db.Account) {
	a.account = account
}

func (a *accountStruct) AccountModel() *db.Account {
	return &a.account
}

func (a *accountStruct) GetAccountInfo() interface{} {

	if len(a.account.Avator) > 0 {
		a.account.Avator = resourceLogic.GenerateToken(a.account.Avator, -1, constants.DaoSuanSessionExpires)
	}

	return paramsUtils.ModelToDict(a.account, field)
}

// 获取关注的人
func (a *accountStruct) GetFollowing() []follow {
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
func (a *accountStruct) GetFollowers() []follow {
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

func (a *accountStruct) GetStars() []dto.ProductList {

	var stars []dto.ProductList
	db.Driver.
		Table("account_star as s, product as p").
		Where("s.account_id = ? and p.id = s.product_id", a.account.Id).
		Select("p.id, p.update_time, p.cover, p.create_time, p.description, p.name, p.status, p.star").
		Order("s.create_time desc").Find(&stars)
	// star产品状态需实时，不适合缓存
	return stars
}

// 检测昵称是否存在
func IsNicknameExists(nickname string, aid ...int) bool {
	if name, err := cache.Dijan.Get(paramsUtils.CacheBuildKey(constants.NicknameModel, nickname)); err == nil && name != nil {
		return true
	}

	var account db.Account
	if len(aid) > 0 {
		if err := db.Driver.Where("nickname = ? and id != ?", nickname, aid[0]).First(&account).Error; err != nil || account.Id == 0 {
			return false
		}
	} else {
		if err := db.Driver.Where("nickname = ?", nickname).First(&account).Error; err != nil || account.Id == 0 {
			return false
		}
	}
	v := hash.RandInt64(240, 240*5)
	cache.Dijan.Set(paramsUtils.CacheBuildKey(constants.NicknameModel, nickname), []byte(nickname), int(v)*60*60)
	return true
}
