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
	"fmt"
	"strconv"
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
	// POST api/v2/user/api_token/create
	ans := make(map[string]interface{})
	args := ToMapStringInterface(t.Args)

	if *t.Method != "POST" {

		response = ErrorReturn(t, 400, "000005", "Wrong Method")
		return response
	}

	p := bluemonday.UGCPolicy()
	tokenName := p.Sanitize(fmt.Sprintf("%v", args["name"]))
	isAdmin := p.Sanitize(fmt.Sprintf("%v", args["isAdmin"]))
	expiration := p.Sanitize(fmt.Sprintf("%v", args["expiration"]))
	comment := p.Sanitize(fmt.Sprintf("%v", args["comment"]))
	readOnly := p.Sanitize(fmt.Sprintf("%v", args["readonly"]))

	if isAdmin == "" {
		isAdmin = "0"
	}
	if expiration == "" {
		expiration = "0"
	}
	if readOnly == "" {
		readOnly = "0"
	}

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

	APIToken := &APITokens{}

	//Create TokenID
	APIToken.TokenId = Hashgen(12)
	APIToken.Created = int(time.Now().Unix())

	//Create Token
	apitkn := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":    t.UID,
		"exipred": expiration,
		"tokenid": APIToken.TokenId,
	})

	APIToken.Token, _ = apitkn.SignedString([]byte(viper.GetString("token.secretKey")))

	//Fill data
	APIToken.UID = *t.UID
	APIToken.Expiration, _ = strconv.Atoi(expiration) //TODO: add conversion from text format to unix time
	APIToken.Status = true

	APIToken.IsAdmin = false
	if isAdmin == "1" {
		isad := false
		if *t.IsAdmin == 1 {
			isad = true
		}
		APIToken.IsAdmin = isad
	}

	redo, _ := strconv.Atoi(readOnly)
	if redo == 1 {
		APIToken.Readonly = true
	} else {
		APIToken.Readonly = false
	}

	APIToken.Comment = comment
	APIToken.TokenName = tokenName

	//Check for table
	if !db.Conn.Migrator().HasTable(&APITokens{}) {

		//Create languages table
		db.Conn.Set("gorm:table_options", "ENGINE=InnoDB;").Migrator().CreateTable(&APITokens{})
	}

	err = db.Conn.Create(&APIToken).Error
	if err != nil {
		response = ErrorReturn(t, 400, "000005", err.Error())
		return response
	}

	ans["tokenid"] = APIToken.TokenId
	ans["created"] = APIToken.Created
	ans["api_token"] = APIToken.Token

	response = Interfacetoresponse(t, ans)
	return response
}
