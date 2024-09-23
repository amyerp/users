package patch

import (
	"fmt"
	"time"
	. "user/grpc_requests"
	. "user/model"

	"github.com/getsentry/sentry-go"
	. "github.com/gogufo/gufo-api-gateway/gufodao"
	"github.com/spf13/viper"

	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/microcosm-cc/bluemonday"
)

func disable2FA(t *pb.Request) (response *pb.Response) {

	args := ToMapStringInterface(t.Args)

	if args["code"] != nil {
		return switchoff2fa(t)
	}

	return initialdisable2fa(t)

}

func initialdisable2fa(t *pb.Request) (response *pb.Response) {
	ans := make(map[string]interface{})

	db, err := ConnectDBv2()
	if err != nil {
		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			SetErrorLog(err.Error())
		}

		return ErrorReturn(t, 500, "000027", err.Error())
	}

	var userExist Users
	db.Conn.Debug().Where(`uid = ?`, *t.UID).First(&userExist)

	if userExist.TFAType == "mail" {

		//sent OTP to email
		otp := Numgen(6)
		lang := "eng"
		if t.Language != nil {
			lang = *t.Language
		}

		go SendOTP(t, userExist.Mail, lang, otp)
		go SendTimeHash(t, otp, userExist.Name, "tfa", userExist.Mail, 300)
	}

	ans["answer"] = "We send OTP to your email.Please check it"
	response = Interfacetoresponse(t, ans)
	return response
}

func switchoff2fa(t *pb.Request) (response *pb.Response) {
	ans := make(map[string]interface{})
	args := ToMapStringInterface(t.Args)
	p := bluemonday.UGCPolicy()

	code := p.Sanitize(fmt.Sprintf("%v", args["code"]))

	db, err := ConnectDBv2()
	if err != nil {
		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			SetErrorLog(err.Error())
		}

		return ErrorReturn(t, 500, "000027", err.Error())
	}

	var userExist Users
	db.Conn.Debug().Where(`uid = ?`, *t.UID).First(&userExist)

	lifetime, _, errstr := CheckTimeHash(t, code, userExist.Name)

	if errstr != "" {
		// return error. user name is exist in db users
		return ErrorReturn(t, 400, "000021", "There is no data")
	}

	// Check for OTP livetime
	ctime := int(time.Now().Unix())

	if ctime > lifetime {
		//Delete OTP
		return ErrorReturn(t, 400, "000022", "OTP has expired")
	}

	go DeleteTimeHash(t, code, userExist.Name)

	db.Conn.Table("users").Where("uid = ?", userExist.UID).Updates(map[string]interface{}{"tfa": int(0)})

	ans["answer"] = "2FA is enabled"
	response = Interfacetoresponse(t, ans)
	return response
}
