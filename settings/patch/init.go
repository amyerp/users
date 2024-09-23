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
	. "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
)

func Init(t *pb.Request) (response *pb.Response) {

	switch *t.ParamID {
	case "dateformat":
		response = setdataresponse(t)
	case "password":
		response = setnewpassword(t)
	case "email":
		response = setnewemail(t)
	case "profile":
		response = updateprofile(t)
	case "enable2fa":
		response = enable2FA(t)
	case "disable2fa":
		response = disable2FA(t)
	default:
		response = ErrorReturn(t, 404, "0000129", "Missing argument")

	}

	return response

}
