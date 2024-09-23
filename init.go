package main

import (
	. "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"

	ad "user/admin"
	tkn "user/api_tokens"
	av "user/avatar"
	pr "user/profile"
	st "user/settings"
	us "user/users"
	. "user/version"
)

func Init(t *pb.Request) (response *pb.Response) {

	if t.UID == nil {
		response = ErrorReturn(t, 401, "000011", "You are not authorised")
		return response
	}

	switch *t.Param {
	case "admin":
		return admincheck(t)
	case "info":
		response = info(t)
	case "health":
		response = health(t)
	case "avatar":
		response = av.Init(t)
	case "users":
		response = us.Init(t)
	case "settings":
		response = st.Init(t)
	case "api_token":
		response = tkn.Init(t)
	case "profile":
		response = pr.Init(t)
	default:
		response = detectbyuser(t)
	}

	/*
		switch *t.Param {
		case "avatar":
			response = av.Init(t)
		case "api_token":
			response = tkn.Init(t)
		case "company":
			response = c.Init(t)
		case "info":
			response = info(t)
		case "get_steps":
			response = Getsteps(t)
		case "hide_steps":
			response = Hidesteps(t)
		case "settings":
			response = set.Init(t)
		case "getusers":
			response = ShowUsers(t)
		case "getuserbyid":
			response = ShowUser(t)
		default:
			ret := fmt.Sprintf("%v: %v", "Missing argument", *t.Param)
			response = ErrorResponseReturn(t, 404, "000121", ret)

		}
	*/
	return response

}

func info(t *pb.Request) (response *pb.Response) {
	ans := make(map[string]interface{})
	ans["pluginname"] = "user"
	ans["version"] = VERSIONPLUGIN
	ans["description"] = ""
	response = Interfacetoresponse(t, ans)
	return response
}

func health(t *pb.Request) (response *pb.Response) {
	ans := make(map[string]interface{})
	ans["health"] = "OK"
	response = Interfacetoresponse(t, ans)
	return response
}

func admincheck(t *pb.Request) (response *pb.Response) {

	if *t.IsAdmin != 1 {
		response = ErrorReturn(t, 401, "000012", "You have no admin rights")
	}

	return ad.Init(t)

}

func detectbyuser(t *pb.Request) (response *pb.Response) {

	if t.ParamID == nil {
		return us.Init(t)
	}

	switch *t.ParamID {
	case "avatar":
		response = av.Init(t)
	case "profile":
		response = pr.Init(t)
	default:
		response = ErrorReturn(t, 406, "000012", "Wrong request")
	}

	return response

}
