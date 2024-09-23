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

	"github.com/getsentry/sentry-go"
	. "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/microcosm-cc/bluemonday"
	"github.com/spf13/viper"
)

// PATCH api/v2/user/api_token/{tokenid}/switch
func UpdateTokenComment(t *pb.Request) (response *pb.Response) {

	ans := make(map[string]interface{})
	args := ToMapStringInterface(t.Args)
	p := bluemonday.UGCPolicy()

	dataid := *t.ParamID
	comment := ""

	if args["comment"] != nil {
		comment = p.Sanitize(fmt.Sprintf("%v", args["comment"]))
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

	curdata := APITokens{}

	rows := db.Conn.Debug().Model(APITokens{}).Where("tokenid = ? AND uid = ?", dataid, *t.UID).First(&curdata)

	if rows.RowsAffected == 0 {
		// return error. user name is exist in db users
		return ErrorReturn(t, 400, "000003", "There is no such token")
	}

	err = db.Conn.Debug().Model(APITokens{}).Where("tokenid = ?", dataid).Update("comment", comment).Error
	if err != nil {
		return ErrorReturn(t, 400, "000005", err.Error())
	}

	//TODO: Record event

	ans["tokenid"] = dataid
	response = Interfacetoresponse(t, ans)
	return response
}
