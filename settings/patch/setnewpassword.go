//////////////////////////////////////////////////////////////////////////////////
// Copyright 2021 Alexey Yanchenko <mail@yanchenko.me>                          //
//                                                                              //
// This file is part of the ERP library.                                        //
//                                                                              //
//  Unauthorized copying of this file, via any media is strictly prohibited     //
//  Proprietary and confidential                                                //
//////////////////////////////////////////////////////////////////////////////////

package patch

import (
	"fmt"
	. "user/model"

	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"golang.org/x/crypto/bcrypt"

	"github.com/getsentry/sentry-go"
	. "github.com/gogufo/gufo-api-gateway/gufodao"
	"github.com/microcosm-cc/bluemonday"
	"github.com/spf13/viper"
)

func setnewpassword(t *pb.Request) (response *pb.Response) {
	// api/user/avatar
	ans := make(map[string]interface{})
	args := ToMapStringInterface(t.Args)

	if args["old_password"] == nil || args["new_password"] == nil {
		return ErrorReturn(t, 406, "000012", "Missing  Important Data")
	}

	p := bluemonday.UGCPolicy()
	oldpassword := p.Sanitize(fmt.Sprintf("%v", args["old_password"]))
	newpassword := p.Sanitize(fmt.Sprintf("%v", args["new_password"]))

	//Check DB and table config
	db, err := ConnectDBv2()
	if err != nil {
		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			SetErrorLog(err.Error())
		}

		response = ErrorReturn(t, 500, "000027", err.Error())
		return response
	}

	//Check old password
	var userExist Users
	db.Conn.Debug().Where(`uid = ?`, *t.UID).First(&userExist)
	if err := bcrypt.CompareHashAndPassword([]byte(userExist.Pass), []byte(oldpassword)); err != nil {
		// Password not matched
		return ErrorReturn(t, 400, "000008", "Password not matched")

	}

	//Hash new new_password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newpassword), 8)
	if err != nil {

		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			SetErrorLog("dbstructure.go: " + err.Error())
		}
	}

	//Update password

	err = db.Conn.Model(&Users{}).Where("uid = ?", *t.UID).Updates(map[string]interface{}{"pass": hashedPassword}).Error
	if err != nil {
		response = ErrorReturn(t, 400, "000005", err.Error())
		return response
	}

	//TODO send notification
	//Create Event

	ans["answer"] = "OK"
	response = Interfacetoresponse(t, ans)
	return response
}
