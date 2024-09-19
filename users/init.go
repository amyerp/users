//////////////////////////////////////////////////////////////////////////////////
// Copyright 2021 Alexey Yanchenko <mail@yanchenko.me>                          //
//                                                                              //
// This file is part of the ERP library.                                        //
//                                                                              //
//  Unauthorized copying of this file, via any media is strictly prohibited     //
//  Proprietary and confidential                                                //
//////////////////////////////////////////////////////////////////////////////////

package users

import (
	pb "github.com/gogufo/gufo-api-gateway/proto/go"

	. "github.com/gogufo/gufo-api-gateway/gufodao"
)

func Init(t *pb.Request) (response *pb.Response) {

	switch *t.Param {
	case "users":
		response = ShowUsers(t)
	default:
		response = ErrorReturn(t, 404, "000014", "Missing Param")
	}
	return response

}
