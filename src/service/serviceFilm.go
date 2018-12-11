package service

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/z774379121/untitled1/src/dao"
	"github.com/z774379121/untitled1/src/xm/common"
	"gopkg.in/mgo.v2/bson"
)

/***
 *    .___________.  ______   .______      .______       _______ .__   __. .___________.    _______  __   __       _______
 *    |           | /  __  \  |   _  \     |   _  \     |   ____||  \ |  | |           |   |   ____||  | |  |     |   ____|
 *    `---|  |----`|  |  |  | |  |_)  |    |  |_)  |    |  |__   |   \|  | `---|  |----`   |  |__   |  | |  |     |  |__
 *        |  |     |  |  |  | |      /     |      /     |   __|  |  . `  |     |  |        |   __|  |  | |  |     |   __|
 *        |  |     |  `--'  | |  |\  \----.|  |\  \----.|  |____ |  |\   |     |  |        |  |     |  | |  `----.|  |____
 *        |__|      \______/  | _| `._____|| _| `._____||_______||__| \__|     |__|        |__|     |__| |_______||_______|
 *
 */
// 上传种子文件,成功返回一个唯一的bid
func UpLoad(context echo.Context) error {
	avatar, err := context.FormFile("avatar")
	if err != nil {
		return err
	}

	// Source
	fmt.Println(avatar.Filename, common.FileSize(avatar.Size))
	src, err := avatar.Open()
	defer src.Close()
	if err != nil {
		return err
	}

	des := make([]byte, avatar.Size)
	n, err2 := src.Read(des)
	if err2 != nil {
		fmt.Println(err2)
		return context.String(http.StatusBadRequest, "读取失败")
	}
	fmt.Println(n)
	// Destination
	daoFile := dao.NewDaoBTFile()
	ok, name := daoFile.UploadBTFile(bson.NewObjectId(), &des)
	if !ok {
		return context.String(http.StatusBadRequest, "插入到数据库失败")
	}

	daoFileContent := dao.NewDaoBTFileContent()
	if daoFileContent.UpdateRealFileName(name, avatar.Filename) {
		return context.JSON(http.StatusOK, map[string]interface{}{
			"msg":  "ok",
			"name": name,
		})
	}
	return context.String(http.StatusBadRequest, "更新名字失败")
}

// 根据上传图片后返回的bson类型文件名下载对应的种子文件
func Download(context echo.Context) error {
	filename := context.Param("filename")
	if !bson.IsObjectIdHex(filename) {
		return context.String(http.StatusForbidden, "invaild filename")
	}

	daoFile := dao.NewDaoBTFile()
	img := daoFile.DownloadBTFile(filename)
	if img == nil {
		return context.String(http.StatusForbidden, "not found")
	}

	daoFileContent := dao.NewDaoBTFileContent()
	imgInfo := daoFileContent.SelectByName(filename)
	filename = filename + ".torrent"
	if imgInfo.RealFileName != "" {
		filename = imgInfo.RealFileName
	}
	fmt.Println(imgInfo.RealFileName)
	dst, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer dst.Close()
	defer os.Remove(filename)

	// Copy
	if _, err = io.Copy(dst, bytes.NewReader(*img)); err != nil {
		return err
	}
	context.Response().Header().Set("Content-Type", "application/octet-stream")
	context.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	return context.File(filename)
}

type FilmStruct struct {
	FilmName   string `form:"film_name"`
	Actor      string `form:"name"`
	Shapeness  int    `form:"shapeness"`
	BTFileName string `form:"bt_file_name"`
	Content    string `form:"content"`
	BTLink     string `form:"bt_link"`
	LocalPath  string `form:"local_path"`
}

func NewFilm(context echo.Context) error {
	f := new(FilmStruct)
	if err := context.Bind(f); err != nil {
		log.Println("表单无效", err)
	}
	token := context.Get("token")

	//daoUser := dao.NewDaoUser()
	//user := daoUser.SelectByAppToken(token.(string))
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
