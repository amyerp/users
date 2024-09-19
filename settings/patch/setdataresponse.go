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

	"user/model"

	pb "github.com/gogufo/gufo-api-gateway/proto/go"

	"github.com/getsentry/sentry-go"
	. "github.com/gogufo/gufo-api-gateway/gufodao"
	"github.com/microcosm-cc/bluemonday"
	"github.com/spf13/viper"
)

func setdataresponse(t *pb.Request) (response *pb.Response) {
	// api/user/avatar
	ans := make(map[string]interface{})
	args := ToMapStringInterface(t.Args)

	csettings := &model.UsersInfo{}

	p := bluemonday.UGCPolicy()
	dateformat := p.Sanitize(fmt.Sprintf("%v", args["dateformat"]))

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

	//Update data to user settings table

	err = db.Conn.Model(&csettings).Where("uid = ?", t.UID).Updates(map[string]interface{}{"dateformat": dateformat}).Error
	if err != nil {
		response = ErrorReturn(t, 400, "000005", err.Error())
		return response
	}

	ans["answer"] = "OK"
	response = Interfacetoresponse(t, ans)
	return response
}
