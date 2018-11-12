package daoImpl

import (
	"dao/baseSeesion"
	"github.com/smallnest/rpcx/log"
	"gopkg.in/mgo.v2/bson"
	"models"
	"models/modelsDefine"
)

type daoSongsImpl struct {
	baseSeesion.BaseSession
}

func NewDaoSongsImpl() *daoSongsImpl {
	obj := &daoSongsImpl{}
	obj.Init(models.Songs{}, models.Collection_Name_Songs)
	return obj
}

func (this *daoSongsImpl) InsertOne(songs *models.Songs) bool {
	return this.BaseSession.InsertModel(songs)
}

func (this *daoSongsImpl) SelectByName(name string) *models.Songs {
	result := &models.Songs{}
	err := this.BaseSession.FindOne(bson.M{modelsDefine.MDSongs_Name: name}, result)
	if err != nil {
		log.Info(err)
	}
	return result
}