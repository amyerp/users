//////////////////////////////////////////////////////////////////////////////////
// Copyright 2021 Alexey Yanchenko <mail@yanchenko.me>                          //
//                                                                              //
// This file is part of the ERP library.                                        //
//                                                                              //
//  Unauthorized copying of this file, via any media is strictly prohibited     //
//  Proprietary and confidential                                                //
//////////////////////////////////////////////////////////////////////////////////

package settings

import (
	pt "user/settings/patch"

	. "github.com/gogufo/gufo-api-gateway/gufodao"

	pb "github.com/gogufo/gufo-api-gateway/proto/go"
)

func Init(t *pb.Request) (response *pb.Response) {

	switch *t.Method {

	case "PATCH":
		response = pt.Init(t)

	default:
		response = ErrorReturn(t, 404, "0000128", "Missing argument")

	}

	return response

}
