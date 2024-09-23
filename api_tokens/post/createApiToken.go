//////////////////////////////////////////////////////////////////////////////////
// Copyright 2021 Alexey Yanchenko <mail@yanchenko.me>                          //
//                                                                              //
// This file is part of the ERP library.                                        //
//                                                                              //
//  Unauthorized copying of this file, via any media is strictly prohibited     //
//  Proprietary and confidential                                                //
//////////////////////////////////////////////////////////////////////////////////

package post

import (
	"encoding/json"
	"fmt"
	"time"

	. "user/model"

	pb "github.com/gogufo/gufo-api-gateway/proto/go"

	"github.com/dgrijalva/jwt-go"
	"github.com/getsentry/sentry-go"
	. "github.com/gogufo/gufo-api-gateway/gufodao"
	"github.com/microcosm-cc/bluemonday"
	"github.com/spf13/viper"
)

func createApiToken(t *pb.Request) (response *pb.Response) {
	ans := make(map[string]interface{})
	args := ToMapStringInterface(t.Args)
	p := bluemonday.UGCPolicy()

	db, err := ConnectDBv2()
	if err != nil {
		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			SetErrorLog(err.Error())
		}

		return ErrorReturn(t, 500, "000027", err.Error())
	}

	newdata := APITokens{}

	JsonArgs, err := json.Marshal(args)
	if err != nil {
		return ErrorReturn(t, 500, "000028", err.Error())
	}

	err = json.Unmarshal(JsonArgs, &newdata)
	if err != nil {
		return ErrorReturn(t, 500, "000028", err.Error())
	}

	tokenid := Hashgen(12)
	newdata.TokenId = tokenid
	newdata.Created = int(time.Now().Unix())

	expiration := p.Sanitize(fmt.Sprintf("%v", args["expiration_string"]))
	if expiration == "" {
		expiration = "0"
	}

	//Create Token
	apitkn := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":    *t.UID,
		"exipred": expiration,
		"tokenid": tokenid,
	})

	newdata.Token, _ = apitkn.SignedString([]byte(viper.GetString("token.secretKey")))
	newdata.UID = *t.UID
	newdata.Expiration = convertTimeToTimestamp(expiration)
	newdata.Status = true

	err = db.Conn.Create(&newdata).Error
	if err != nil {
		response = ErrorReturn(t, 400, "000005", err.Error())
		return response
	}

	ans["tokenid"] = newdata.TokenId
	ans["created"] = newdata.Created
	ans["api_token"] = newdata.Token

	response = Interfacetoresponse(t, ans)
	return response
}

func convertTimeToTimestamp(date string) int {
	//Check the documentation on Go for the const variables!
	//They need to be exactly as they are shown in the documentation to be read correctly!
	format := "2006-01-02"

	t, err := time.Parse(format, date)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(t.Unix())
		return int(t.Unix())
	}
	return 0
}
