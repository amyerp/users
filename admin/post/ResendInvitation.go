package post

import (
	"fmt"
	"time"
	. "user/grpc_requests"
	. "user/model"

	"github.com/getsentry/sentry-go"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/microcosm-cc/bluemonday"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"

	. "github.com/gogufo/gufo-api-gateway/gufodao"
)

func ResendInvitation(t *pb.Request) (response *pb.Response) {
	ans := make(map[string]interface{})
	args := ToMapStringInterface(t.Args)
	p := bluemonday.UGCPolicy()

	domain := viper.GetString("server.domain")

	//Check deoes Username, and email provided
	if args["uid"] == nil {
		return ErrorReturn(t, 404, "000012", "Please provide User ID")
	}

	uid := p.Sanitize(fmt.Sprintf("%v", args["uid"]))

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
	db.Conn.Debug().Where("uid = ?", uid).First(&curdata)

	curtime := int(time.Now().Unix())
	sandmailtime := curdata.Mailsent
	if curdata.Access != 0 {
		return ErrorReturn(t, 406, "000012", "User already confirmed his account")
	}

	dif := curtime - sandmailtime
	waittime := 300
	if dif < waittime {
		return ErrorReturn(t, 406, "000012", "You can resend email again in 5 minutes")
	}

	password := Hashgen(14)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {

		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			SetErrorLog("dbstructure.go: " + err.Error())
		}
	}
	newdata := Users{}
	newdata.Pass = string(hashedPassword)
	newdata.Mailsent = int(time.Now().Unix())
	if args["mail"] != nil {
		umail := p.Sanitize(fmt.Sprintf("%v", args["mail"]))
		newdata.Mail = umail
	}

	err = db.Conn.Where("uid = ?", uid).Updates(&newdata).Error
	if err != nil {
		return ErrorReturn(t, 400, "000005", err.Error())
	}

	//We should send email to user with new password. And generate random password if pasword not provided
	title := "Invitation to Amy ERP"
	message := []string{}
	msga := "You received this email because somebody create an Account for you."
	message = append(message, msga)
	msga = "Please use nxt credentials for your first login:"
	message = append(message, msga)
	msga = fmt.Sprintf("host: %s", domain)
	message = append(message, msga)
	msga = fmt.Sprintf("login: %s", curdata.Name)
	message = append(message, msga)
	msga = fmt.Sprintf("password: <code>%s</code>", password)
	message = append(message, msga)
	msga = "Do not forget to change your password after first login!"
	message = append(message, msga)

	go SendNotification(t, title, message, "invitation", uid)

	ans["uid"] = uid
	return Interfacetoresponse(t, ans)

}
