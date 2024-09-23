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
	switch *t.ParamIDD {
	case "switch":
		response = SwithTokenStatus(t)
	case "comment":
		response = UpdateTokenComment(t)
	default:
		response = ErrorReturn(t, 404, "000014", "Missing Param")
	}

	return response

}
