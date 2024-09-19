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
	"encoding/json"
	"user/model"

	pb "github.com/gogufo/gufo-api-gateway/proto/go"

	"github.com/getsentry/sentry-go"
	. "github.com/gogufo/gufo-api-gateway/gufodao"
	"github.com/spf13/viper"
)

func updateprofile(t *pb.Request) (response *pb.Response) {
	// api/user/avatar
	ans := make(map[string]interface{})
	args := ToMapStringInterface(t.Args)

	//TODO: check that not alowed to change UID, Role and PersonID

	csettings := &model.UsersInfo{}

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

	JsonArgs, err := json.Marshal(args)
	if err != nil {
		return ErrorReturn(t, 500, "000028", err.Error())
	}

	err = json.Unmarshal(JsonArgs, &csettings)
	if err != nil {
		return ErrorReturn(t, 500, "000028", err.Error())
	}

	err = db.Conn.Where("uuid = ?", *t.UID).Updates(&csettings).Error
	if err != nil {
		return ErrorReturn(t, 400, "000005", err.Error())
	}

	ans["answer"] = "OK"
	response = Interfacetoresponse(t, ans)
	return response
}
