package avatar

import (
	. "user/functions"
	. "user/model"

	"github.com/gabriel-vasile/mimetype"
	"github.com/getsentry/sentry-go"
	. "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/spf13/viper"
)

func GetAvatar(t *pb.Request) (response *pb.Response) {
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

	if *t.Param != "avatar" {
		//means that we need avatar of another users
		userid = *t.Param

		//Check user
		var count int64
		db.Conn.Debug().Model(Users{}).Where("uid = ?", userid).Count(&count)
		if count == 0 {
			return ErrorReturn(t, 404, "000023", "User Not Found")
		}
	}

	user := &UsersInfo{}
	rows := db.Conn.Debug().Model(&user).Where("uid = ?", userid).First(&user)
	if rows.RowsAffected == 0 {
		return ErrorReturn(t, 406, "000027", "User not found")
	}

	link := user.AvatarID
	if link == "" {
		link = GenUserAvatar(userid, t)
	}

	mime, err := mimetype.DetectFile(link)
	if err != nil {

		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			SetErrorLog("Detect file error: " + err.Error())
		}

		return ErrorReturn(t, 400, "0000035", err.Error())

	}

	ans["file"] = link
	ans["filetype"] = mime.String()
	ans["filename"] = "avatar.png"

	response = Interfacetoresponse(t, ans)

	return response
}
