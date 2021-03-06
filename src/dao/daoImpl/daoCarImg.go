package daoImpl

import (
	"github.com/z774379121/untitled1/src/dao/baseSession"
	"github.com/z774379121/untitled1/src/models"
	"gopkg.in/mgo.v2/bson"
)

type daoCarImgImpl struct {
	baseSession.BaseSession
}

func NewDaoCarImgImpl() *daoCarImgImpl {
	obj := &daoCarImgImpl{}
	obj.Init(models.CardCouponImg{}, models.COLLECTION_NAME_CardCouponImg)
	return obj
}

func (this *daoCarImgImpl) UploadImg(songId bson.ObjectId, filedata *[]byte) (bool, string) {
	filename := bson.NewObjectId().Hex()
	this.CreateFile(filename, *filedata)
	return true, filename
}
func (this *daoCarImgImpl) DownloadImg(filename string) *[]byte {
	data := this.ReadFile(filename)
	return data
}
