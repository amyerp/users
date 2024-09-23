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
	"time"
	. "user/model"

	. "user/grpc_requests"

	"github.com/getsentry/sentry-go"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"

	. "github.com/gogufo/gufo-api-gateway/gufodao"
	"github.com/microcosm-cc/bluemonday"
)

func setnewemail(t *pb.Request) (response *pb.Response) {

	args := ToMapStringInterface(t.Args)

	if args["code"] != nil {
		return updatenewemail(t)
	}

	return initialsetnewemail(t)
}

func initialsetnewemail(t *pb.Request) (response *pb.Response) {
	// api/user/avatar
	ans := make(map[string]interface{})
	args := ToMapStringInterface(t.Args)

	ischange := viper.GetBool("settings.change_email")

	if !ischange {
		return ErrorReturn(t, 403, "000010", "You are not allowed to change email. Please contact your Administrator with such request")
	}

	if args["email"] == nil || args["password"] == nil {
		return ErrorReturn(t, 406, "000012", "Missing  Important Data")
	}

	p := bluemonday.UGCPolicy()
	email := p.Sanitize(fmt.Sprintf("%v", args["email"]))
	password := p.Sanitize(fmt.Sprintf("%v", args["password"]))

	//Check does suck email is exist in DB
	db, err := ConnectDBv2()
	if err != nil {
		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			SetErrorLog(err.Error())
		}
	}

	//Check old password
	var userExist Users
	db.Conn.Debug().Where(`uid = ?`, *t.UID).First(&userExist)
	if err := bcrypt.CompareHashAndPassword([]byte(userExist.Pass), []byte(password)); err != nil {
		// Password not matched
		return ErrorReturn(t, 400, "000008", "Password not matched")

	}

	var count int64
	db.Conn.Debug().Model(Users{}).Where("mail = ?", email).Count(&count)

	if count != 0 {
		return ErrorReturn(t, 406, "000012", "Such Email is already exist")
	}

	//generate code
	otp := Numgen(6)

	lang := "eng"
	if t.Language != nil {
		lang = *t.Language
	}

	//send code to email
	go SendOTP(t, email, lang, otp)
	go SendTimeHash(t, otp, userExist.Name, "tfa", email, 300)

	ans["response"] = "100201"
	response = Interfacetoresponse(t, ans)
	return response

}

func updatenewemail(t *pb.Request) (response *pb.Response) {
	args := ToMapStringInterface(t.Args)
	p := bluemonday.UGCPolicy()
	ans := make(map[string]interface{})

	if args["email"] == nil {
		return ErrorReturn(t, 406, "000012", "Missing  email")
	}

	code := p.Sanitize(fmt.Sprintf("%v", args["code"]))
	email := p.Sanitize(fmt.Sprintf("%v", args["email"]))

	lifetime, _, errstr := CheckTimeHash(t, code, email)

	if errstr != "" {
		// return error. user name is exist in db users
		return ErrorReturn(t, 404, "000021", "Code not find")
	}

	// Check for OTP livetime
	ctime := int(time.Now().Unix())
	SetErrorLog(fmt.Sprintf("ctime: %v", ctime))
	SetErrorLog(fmt.Sprintf("lifetime: %v", lifetime))

	if lifetime < ctime {
		//Delete OTP
		go DeleteTimeHash(t, code, email)
		return ErrorReturn(t, 400, "000022", "Code has expired")
	}

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

	//Update email in users table
	err = db.Conn.Model(&Users{}).Where("uid = ?", *t.UID).Updates(map[string]interface{}{"mail": email, "mailconfirmed": int(time.Now().Unix())}).Error
	if err != nil {
		response = ErrorReturn(t, 400, "000005", err.Error())
		return response
	}

	//TODO: check does user has PersonID and check for such email in email list. if not exist, add new email to email list

	ans["answer"] = "OK"
	response = Interfacetoresponse(t, ans)
	return response

}
