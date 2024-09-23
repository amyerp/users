package users

import (
	. "user/model"

	"github.com/getsentry/sentry-go"
	. "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/spf13/viper"
)

func ShowUser(t *pb.Request) (response *pb.Response) {

	ans := make(map[string]interface{})

	//p := bluemonday.UGCPolicy()

	user := Users{}
	userinfo := UsersInfo{}

	uid := *t.Param
	if uid == "" {
		return ErrorReturn(t, 404, "000013", "Missing User ID")
	}

	db, err := ConnectDBv2()
	if err != nil {
		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			SetErrorLog(err.Error())
		}
		return ErrorReturn(t, 500, "000027", err.Error())
	}

	var count int64

	db.Conn.Debug().Model(Users{}).Where("uid = ?", uid).Count(&count)
	if count == 0 {
		return ErrorReturn(t, 404, "000013", "User not Found")
	}

	db.Conn.Debug().Where("uid = ?", uid).First(&user)
	db.Conn.Debug().Where("uid = ?", uid).First(&userinfo)
	user.Pass = ""

	if *t.IsAdmin != 1 {

		user.Mailsent = 0
		user.Mailconfirmed = 0
		user.IP = ""
		user.Access = 0
		user.Completed = false
		user.TFA = false
		user.TFAType = ""

	}

	curuser := UserResponse{}
	curuser.Users = &user
	curuser.UsersInfo = &userinfo

	ans["user"] = curuser
	response = Interfacetoresponse(t, ans)
	return response

}
