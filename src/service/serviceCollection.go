package service

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/z774379121/untitled1/src/dao"
	"github.com/z774379121/untitled1/src/logger"
	"github.com/z774379121/untitled1/src/models"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"strings"
	"time"
)

func NewFilm(context echo.Context) error {
	f := new(FNewFilm)
	if err := context.Bind(f); err != nil {
		log.Println("表单无效", err)
	}
	token := context.Get("token")
	publish := f.Publisher
	if f.IsPrimitive {
		publish = f.Publisher1
	}
	parse, e := time.Parse("2006-01", f.Date)
	if e != nil {
		return e
	}

	daoUser := dao.NewDaoUser()
	daoFilm := dao.NewDaoFilm()
	exsistFilm := daoFilm.FindByCode(f.Code)
	if exsistFilm != nil {
		return context.String(http.StatusOK, "已经存在该作品， 谢谢你的参与")
	}
	user := daoUser.SelectByAppToken(token.(string))
	newFilm := models.NewFilm()
	newFilm.Name = f.Code
	newFilm.ActorRef.Id = bson.ObjectIdHex(f.Actor)
	newFilm.IsPrimitive = f.IsPrimitive
	newFilm.Publisher = models.Publishes(publish)
	newFilm.Year = parse.Year()
	newFilm.FriendShip = strings.Split(f.Serice, ",")
	inserted := daoFilm.InsertOne(newFilm)
	if !inserted {
		return context.String(http.StatusBadRequest, "bad")
	}
	fmt.Println(user.Email)

	////daoClt := dao.NewColletion()
	//daoFile := dao.NewDaoBTFileContent()
	//daoFile.SelectByName(f.BTFileName)
	//
	//newObj := models.NewColletion()
	////newObj.FilmRef.Id_ = f.FilmName
	//newObj.Content = f.Content
	//
	//newObj.BTLink = f.BTLink
	//newObj.LocalPath = f.LocalPath
	//newObj.Shapeness = models.Shapeness(f.Shapeness)
	//newObj.UserRef.Id = user.Id_
	fmt.Println(token)
	fmt.Println(f)
	return context.String(http.StatusOK, "ok")
}

func GetColletions(ctx echo.Context) error {
	token := ctx.Get("token")
	daoUser := dao.NewDaoUser()
	user := daoUser.SelectByAppToken(token.(string))
	if user == nil {
		return ctx.String(http.StatusBadRequest, "bad uid")
	}

	daoCollection := dao.NewDaoCollection()
	collections := daoCollection.SelectByUid(user.Id_)
	if collections != nil {
		return ctx.String(http.StatusBadRequest, "无效查询")
	}

	return ctx.JSON(http.StatusOK, collections)
}

func GetCollectionUnderTag(ctx echo.Context) error {
	token := ctx.Get("token")
	daoUser := dao.NewDaoUser()
	user := daoUser.SelectByAppToken(token.(string))
	if user == nil {
		return ctx.String(http.StatusBadRequest, "bad uid")
	}
	tag := ctx.QueryParam("tag")
	if !bson.IsObjectIdHex(tag) {
		return ctx.String(http.StatusBadRequest, "bad tagId")
	}
	daoCollection := dao.NewDaoCollection()
	collections := daoCollection.SelectByTagIdsAndUid(user.Id_, []bson.ObjectId{bson.ObjectIdHex(tag)})
	return ctx.JSON(http.StatusOK, collections)
}

type RemoveCTagOpt struct {
	tagId string `form:"tag_id"`
	colId string `form:"col_id"`
}

func RemoveCollectionTag(ctx echo.Context) error {
	token := ctx.Get("token")
	daoUser := dao.NewDaoUser()
	user := daoUser.SelectByAppToken(token.(string))
	if user == nil {
		return ctx.String(http.StatusBadRequest, "bad uid")
	}
	opt := new(RemoveCTagOpt)
	if err := ctx.Bind(opt); err != nil {
		logger.Sugar.Error(err)
	}
	if !bson.IsObjectIdHex(opt.tagId) || !bson.IsObjectIdHex(opt.colId) {
		return echo.ErrBadRequest
	}
	daoCol := dao.NewDaoCollection()
	if done := daoCol.RemoveTagByTagIdAndUidCollectionId(bson.ObjectIdHex(opt.tagId), user.Id_, bson.ObjectIdHex(opt.colId)); !done {
		logger.Sugar.Infow("删除失败",
			"cid", opt.colId,
			"tid", opt.tagId,
			"uid", user.Id_)
		return echo.NewHTTPError(http.StatusBadRequest, "删除失败")
	}
	return ctx.String(http.StatusOK, "移除标签成功")
}
