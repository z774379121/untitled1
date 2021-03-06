package models

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"time"

	"github.com/z774379121/untitled1/src/xm/common"
	"golang.org/x/crypto/pbkdf2"
	"gopkg.in/mgo.v2/bson"
)

const COLLECTION_NAME_User string = "users"

const (
	HasSuperAuAuthorization_Login        = 1 << iota //表示拥有超级管理员的登录权限
	HasSuperAuAuthorization_GetUserData              //表示拥有获取用户数据的权限
	HasSuperAuAuthorization_EditUserData             //表示拥有修改用户数据的权限
)

var DefaultAuthor = HasSuperAuAuthorization_GetUserData | HasSuperAuAuthorization_Login

type User struct {
	Id_           bson.ObjectId `bson:"_id"`
	Email         string        `bson:"email" ` //邮箱
	Username      string        `bson:"username"`
	IsConfirmed   bool          `bson:"is_confirmed"`
	Password      string        `bson:"password"` //密码，以md5形式储存
	Salt          string        `bson:"salt"`
	Session       string        `bson:"session"`
	Token         string        `bson:"token"` //登录令牌
	Authorization int
	CreateTime    time.Time `bson:"create_time"`

	testGenUserToken func() (session, token string) `bson:",omitempty"` //测试注入用，处理用于生成测试的token
}

var testUser *User

func NewUser() *User {
	if testUser != nil {
		return testUser
	}
	obj := &User{}
	obj.Id_ = bson.NewObjectId()
	obj.CreateTime = time.Now()
	obj.Authorization = DefaultAuthor
	return obj
}

func (this *User) GenUserToken() {
	if this.testGenUserToken == nil {
		uuid, _ := common.GenUUID32()
		this.Session = uuid
		this.Token = common.GetMD5(this.Session + this.Id_.Hex())
	} else {
		this.Session, this.Token = this.testGenUserToken()
	}
}

func (this *User) GetAppToken() string {
	return this.Session + this.Token
}

func (this *User) GenSalt() {
	this.Salt = common.RandString(10)
}

func (u *User) EncodePasswd() {
	newPasswd := pbkdf2.Key([]byte(u.Password), []byte(u.Salt), 10000, 50, sha256.New)
	u.Password = fmt.Sprintf("%x", newPasswd)
}

// ValidatePassword checks if given password matches the one belongs to the user.
func (u *User) ValidatePassword(passwd string) bool {
	newUser := &User{Password: passwd, Salt: u.Salt}
	newUser.EncodePasswd()
	return subtle.ConstantTimeCompare([]byte(u.Password), []byte(newUser.Password)) == 1
}
