package users

import (
	"fmt"
	"strconv"
	. "user/model"

	"github.com/getsentry/sentry-go"
	. "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/spf13/viper"
)

func ShowUsers(t *pb.Request) (response *pb.Response) {

	ans := make(map[string]interface{})
	args := ToMapStringInterface(t.Args)
	//p := bluemonday.UGCPolicy()

	user := []Users{}

	db, err := ConnectDBv2()
	if err != nil {
		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			SetErrorLog(err.Error())
		}
		return ErrorReturn(t, 500, "000027", err.Error())
	}

	offset := 0
	limit := 25

	if args["offset"] != nil {
		offset, _ = strconv.Atoi(fmt.Sprintf("%v", args["offset"]))
	}

	if args["limit"] != nil {
		limit, _ = strconv.Atoi(fmt.Sprintf("%v", args["limit"]))
	}

	var count int64
	if args["status"] != nil {
		status := fmt.Sprintf("%v", args["status"])
		db.Conn.Debug().Model(Users{}).Where("status = ?", status).Count(&count)
		db.Conn.Debug().Where("status = ?", status).Limit(limit).Offset(offset).Find(&user)
	} else {
		db.Conn.Debug().Model(Users{}).Count(&count)
		db.Conn.Debug().Limit(limit).Offset(offset).Find(&user)
	}

	allusers := []UserResponse{}

	for i := 0; i < len(user); i++ {
		user[i].Pass = ""
		user[i].Mailsent = 0
		user[i].Mailconfirmed = 0
		user[i].IP = ""
		user[i].Access = 0
		user[i].Completed = false
		user[i].TFA = false
		user[i].TFAType = ""

		userinfo := UsersInfo{}

		db.Conn.Debug().Where("uid = ?", user[i].UID).First(&userinfo)
		userinfo.AvatarID = ""
		userinfo.UID = ""

		curuser := UserResponse{}
		curuser.Users = &user[i]
		curuser.UsersInfo = &userinfo

		allusers = append(allusers, curuser)
	}

	ans["userscount"] = count
	ans["users"] = allusers
	response = Interfacetoresponse(t, ans)
	return response

}
