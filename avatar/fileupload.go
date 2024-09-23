//////////////////////////////////////////////////////////////////////////////////
// Copyright 2021 Alexey Yanchenko <mail@yanchenko.me>                          //
//                                                                              //
// This file is part of the ERP library.                                        //
//                                                                              //
//  Unauthorized copying of this file, via any media is strictly prohibited     //
//  Proprietary and confidential                                                //
//////////////////////////////////////////////////////////////////////////////////

package avatar

import (
	"os"
	"path/filepath"
	"user/model"

	"github.com/getsentry/sentry-go"
	. "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/spf13/viper"
)

func fileUpload(t *pb.Request) (response *pb.Response) {

	if *t.Param != "avatar" {
		return ErrorReturn(t, 403, "000023", "Operation not permited")
	}

	ans := make(map[string]interface{})

	//1. Check for company ID
	//Check DB and table config
	db, err := ConnectDBv2()
	if err != nil {
		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			SetErrorLog(err.Error())
		}

		return ErrorReturn(t, 500, "000027", err.Error())
	}

	fileid := Hashgen(12)
	filename := *t.Filename
	extension := filepath.Ext(filename)

	var pwd string = viper.GetString("server.filedir")

	pathfiles := filepath.Join(pwd, "users")
	//Create dir output using above code
	if _, err := os.Stat(pathfiles); os.IsNotExist(err) {
		os.Mkdir(pathfiles, 0755)
	}

	pathfiles = filepath.Join(pathfiles, *t.UID)
	//Create dir output using above code
	if _, err := os.Stat(pathfiles); os.IsNotExist(err) {
		os.Mkdir(pathfiles, 0755)
	}

	pathfiles = filepath.Join(pathfiles, "avatar")
	//Create dir output using above code
	if _, err := os.Stat(pathfiles); os.IsNotExist(err) {
		os.Mkdir(pathfiles, 0755)
	}

	filelink := pathfiles + "/" + fileid + extension

	f, err := os.OpenFile(filelink, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {

		return ErrorReturn(t, 400, "000026", err.Error())
	}
	defer f.Close()
	f.Write(t.File)

	userinfo := &model.UsersInfo{}
	userinfo.AvatarID = filelink
	db.Conn.Where("uid = ?", *t.UID).Updates(&userinfo)

	ans["status"] = "OK"
	return Interfacetoresponse(t, ans)

}
