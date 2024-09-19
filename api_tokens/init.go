//////////////////////////////////////////////////////////////////////////////////
// Copyright 2021 Alexey Yanchenko <mail@yanchenko.me>                          //
//                                                                              //
// This file is part of the ERP library.                                        //
//                                                                              //
//  Unauthorized copying of this file, via any media is strictly prohibited     //
//  Proprietary and confidential                                                //
//////////////////////////////////////////////////////////////////////////////////

package api_tokens

import (
	. "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"

	dl "user/api_tokens/delete"
	gt "user/api_tokens/get"
	pt "user/api_tokens/patch"
	ps "user/api_tokens/post"
)

func Init(t *pb.Request) (response *pb.Response) {

	switch *t.Method {
	case "POST":
		response = ps.Init(t)
	case "GET":
		response = gt.ShowApiToken(t)
	case "PATCH":
		response = pt.UpdateApiToken(t)
	case "DELETE":
		response = dl.DeleteApiToken(t)
	default:
		response = ErrorReturn(t, 404, "0000128", "Missing argument")
	}

	return response

}
