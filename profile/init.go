package profile

import (
	//. "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
)

func Init(t *pb.Request) (response *pb.Response) {

	return ShowUser(t)
	/*
	   switch *t.Param {
	   case "profile":

	   	response = ShowUser(t)

	   default:

	   		response = ErrorReturn(t, 404, "000014", "Missing Param")
	   	}

	   return response
	*/
}
