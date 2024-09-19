package profile

import (
	. "user/model"

	"github.com/getsentry/sentry-go"
	. "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/spf13/viper"
)

func ShowUser(t *pb.Request) (response *pb.Response) {

	ans := make(map[string]interface{})

	userid := *t.UID

	db, err := ConnectDBv2()
	if err != nil {
		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			SetErrorLog(err.Error())
		}
		return ErrorReturn(t, 500, "000027", err.Error())
	}

	if *t.Param != "profile" {
		//means that we need avatar of another users
		userid = *t.Param

		//Check user
		var count int64
		db.Conn.Debug().Model(Users{}).Where("uid = ?", userid).Count(&count)
		if count == 0 {
			return ErrorReturn(t, 404, "000023", "User Not Found")
		}
	}

	userdata := UserResponse{}
	user := Users{}

	db.Conn.Debug().Where("uid = ?", userid).First(&user)

	user.Pass = ""

	if *t.UID != userid {
		user.Mailsent = 0
		user.Mailconfirmed = 0
		user.Login = 0
		user.IP = ""
		user.Access = 0
		user.Completed = false
		user.TFA = false
		user.TFAType = ""
	}

	userinfo := UsersInfo{}

	db.Conn.Debug().Where("uid = ?", userid).First(&userinfo)
	userinfo.AvatarID = ""
	userinfo.UID = ""

	userdata.Users = &user
	userdata.UsersInfo = &userinfo

	ans["profile"] = userdata
	response = Interfacetoresponse(t, ans)
	return response

}
