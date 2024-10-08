//////////////////////////////////////////////////////////////////////////////////
// Copyright 2021 Alexey Yanchenko <mail@yanchenko.me>                          //
//                                                                              //
// This file is part of the ERP library.                                        //
//                                                                              //
//  Unauthorized copying of this file, via any media is strictly prohibited     //
//  Proprietary and confidential                                                //
//////////////////////////////////////////////////////////////////////////////////

package delete

import (
	. "user/model"

	. "github.com/gogufo/gufo-api-gateway/gufodao"

	"github.com/getsentry/sentry-go"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/spf13/viper"
)

func DeleteApiToken(t *pb.Request) (response *pb.Response) {

	ans := make(map[string]interface{})

	dataid := *t.ParamID

	db, err := ConnectDBv2()
	if err != nil {
		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			SetErrorLog(err.Error())
		}

		return ErrorReturn(t, 500, "000027", err.Error())
	}

	curdata := APITokens{}

	rows := db.Conn.Debug().Model(APITokens{}).Where("tokenid = ? AND uid = ?", dataid, *t.UID).First(&curdata)

	if rows.RowsAffected == 0 {
		// return error. user name is exist in db users
		return ErrorReturn(t, 400, "000003", "There is no such token")
	}

	vc := APITokens{}
	db.Conn.Where("tokenid = ?", dataid).Delete(&vc)

	//TODO: Record event

	ans["answer"] = "OK"
	response = Interfacetoresponse(t, ans)
	return response
}
