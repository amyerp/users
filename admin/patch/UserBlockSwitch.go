package patch

import (
	"fmt"
	. "user/model"

	. "github.com/gogufo/gufo-api-gateway/gufodao"

	"github.com/getsentry/sentry-go"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/microcosm-cc/bluemonday"
	"github.com/spf13/viper"
)

func UserBlockSwitch(t *pb.Request) (response *pb.Response) {

	ans := make(map[string]interface{})
	args := ToMapStringInterface(t.Args)
	p := bluemonday.UGCPolicy()

	if args["uid"] == nil {
		return ErrorReturn(t, 406, "000012", "Missing  User ID")
	}

	dataid := p.Sanitize(fmt.Sprintf("%v", args["uid"]))

	db, err := ConnectDBv2()
	if err != nil {
		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			SetErrorLog(err.Error())
		}

		return ErrorReturn(t, 500, "000027", err.Error())
	}

	curdata := Users{}

	db.Conn.Debug().Model(Users{}).Where("uid = ?", dataid).First(&curdata)
	/*
		data := Users{}

		if curdata.Status {
			data.Status = false
		} else {
			data.Status = true
		}
	*/
	err = db.Conn.Debug().Model(Users{}).Where("uid = ?", dataid).Update("status", !curdata.Status).Error
	if err != nil {
		return ErrorReturn(t, 400, "000005", err.Error())
	}

	//TODO: Record event

	ans["uuid"] = dataid
	response = Interfacetoresponse(t, ans)
	return response
}
