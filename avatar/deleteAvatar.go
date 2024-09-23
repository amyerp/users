package avatar

import (
	. "user/functions"

	. "github.com/gogufo/gufo-api-gateway/gufodao"

	pb "github.com/gogufo/gufo-api-gateway/proto/go"
)

func deleteAvatar(t *pb.Request) (response *pb.Response) {
	ans := make(map[string]interface{})

	if *t.Param != "avatar" {
		return ErrorReturn(t, 403, "000023", "Operation not permited")
	}

	GenUserAvatar(*t.UID, t)

	ans["status"] = "OK"

	return Interfacetoresponse(t, ans)

}
