//////////////////////////////////////////////////////////////////////////////////
// Copyright 2021 Alexey Yanchenko <mail@yanchenko.me>                          //
//                                                                              //
// This file is part of the ERP library.                                        //
//                                                                              //
//  Unauthorized copying of this file, via any media is strictly prohibited     //
//  Proprietary and confidential                                                //
//////////////////////////////////////////////////////////////////////////////////

package get

import (
	"fmt"
	"strconv"
	"strings"

	. "user/model"

	pb "github.com/gogufo/gufo-api-gateway/proto/go"

	"github.com/getsentry/sentry-go"
	. "github.com/gogufo/gufo-api-gateway/gufodao"
	"github.com/microcosm-cc/bluemonday"
	"github.com/spf13/viper"
)

func ShowApiToken(t *pb.Request) (response *pb.Response) {
	// GET api/v2/user/api_token/show
	// GET api/v2/user/api_token/show/{tokenid}/
	// api/v2/{t.Module}/{t.Param}/{t.ParamID}

	path := *t.Path
	patharray := strings.Split(path, "/")
	pathlenth := len(patharray)

	if pathlenth >= 7 {
		response = showtoken(t)

	} else {
		response = showtokens(t)
	}

	return response

}

func showtoken(t *pb.Request) (response *pb.Response) {
	ans := make(map[string]interface{})

	path := *t.Path
	patharray := strings.Split(path, "/")

	p := bluemonday.UGCPolicy()

	tokenid := p.Sanitize(patharray[6])

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

	apitoken := APITokens{}
	db.Conn.Where(`tokenid = ? AND uid = ?`, tokenid, t.UID).Find(&apitoken)

	ans["api_token"] = apitoken

	response = Interfacetoresponse(t, ans)
	return response
}

func showtokens(t *pb.Request) (response *pb.Response) {
	ans := make(map[string]interface{})
	args := ToMapStringInterface(t.Args)

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

	offset := 0
	limit := 25

	if args["offset"] != nil {
		offset, _ = strconv.Atoi(fmt.Sprintf("%v", args["offset"]))
	}

	if args["limit"] != nil {
		limit, _ = strconv.Atoi(fmt.Sprintf("%v", args["limit"]))
	}

	var count int64
	db.Conn.Model(APITokens{}).Where(`uid = ?`, t.UID).Count(&count)

	alltokens := []APITokens{}
	db.Conn.Where(`uid = ?`, t.UID).Limit(limit).Offset(offset).Find(&alltokens)

	ans["api_tokens"] = alltokens
	ans["tokenscount"] = count

	response = Interfacetoresponse(t, ans)
	return response
}
