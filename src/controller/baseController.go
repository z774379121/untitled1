package controller

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/z774379121/untitled1/src/dao"
	"github.com/z774379121/untitled1/src/models"
	"net/http"
	"time"
)

const (
	Key_cookie_UserToken  = "user_token"
	Key_cookie_Token_Type = "token_type"
	Key_cookie_Time       = "_time"
)

type BaseController struct {
	UserToken string
	TokenType string
	C         echo.Context
}

func NewBaseController(ctx echo.Context) *BaseController {
	return &BaseController{
		C: ctx,
	}
}

func (this *BaseController) SetToken(tokenType, userToken string) {
	this.TokenType = tokenType
	this.UserToken = userToken
}

func (this *BaseController) GetUserToken() string {
	cookie, err := this.C.Cookie(Key_cookie_UserToken)
	if err != nil {
		fmt.Println("不存在")
		return ""
	}
	return cookie.Value
}

func (this *BaseController) GetTokenType() string {
	cookie, err := this.C.Cookie(Key_cookie_Token_Type)
	if err != nil || cookie == nil {
		fmt.Println("不存在")
		return ""
	}
	return cookie.Value
}

func (this *BaseController) SetCookies() {
	this.C.SetCookie(&http.Cookie{Name: Key_cookie_UserToken, Value: this.UserToken, Expires: time.Now().Add(time.Hour * 2), Path: "/"})
	this.C.SetCookie(&http.Cookie{Name: Key_cookie_Time, Value: time.Now().Format("20060102150405"), Expires: time.Now().Add(time.Hour * 2), Path: "/"})
	this.C.SetCookie(&http.Cookie{Name: Key_cookie_Token_Type, Value: this.TokenType, Expires: time.Now().Add(time.Hour * 2), Path: "/"})
}

func (this *BaseController) ClearCookies() {
	this.C.SetCookie(&http.Cookie{Name: Key_cookie_UserToken, Value: "", Expires: time.Now(), Path: "/"})
	this.C.SetCookie(&http.Cookie{Name: Key_cookie_Time, Value: "", Expires: time.Now(), Path: "/"})
	this.C.SetCookie(&http.Cookie{Name: Key_cookie_Token_Type, Value: "", Expires: time.Now(), Path: "/"})
}

func GetUser(token string) *models.User {
	daoUser := dao.NewDaoUser()
	return daoUser.SelectByAppToken(token)
}
