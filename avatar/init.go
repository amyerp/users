//////////////////////////////////////////////////////////////////////////////////
// Copyright 2021 Alexey Yanchenko <mail@yanchenko.me>                          //
//                                                                              //
// This file is part of the ERP library.                                        //
//                                                                              //
//  Unauthorized copying of this file, via any media is strictly prohibited     //
//  Proprietary and confidential                                                //
//////////////////////////////////////////////////////////////////////////////////

package avatar

import (
	"fmt"

	pb "github.com/gogufo/gufo-api-gateway/proto/go"

	. "github.com/gogufo/gufo-api-gateway/gufodao"
)

func Init(t *pb.Request) (response *pb.Response) {
	SetLog("Avatar module")

	switch *t.Method {
	case "GET":
		// show compnay
		response = GetAvatar(t) // No such param

	case "PUT":
		// add or edit comapny
		response = fileUpload(t)

	case "DELETE":
		// delete company
		//	ans, t = delete(t)
		response = deleteAvatar(t)
	default:
		ret := fmt.Sprintf("%v: %v", "Missing argument", *t.Param)
		response = ErrorReturn(t, 404, "000122", ret)

	}

	return response

}
