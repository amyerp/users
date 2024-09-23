package post

import (
	"encoding/json"
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

func CreateUser(t *pb.Request) (response *pb.Response) {
	ans := make(map[string]interface{})
	args := ToMapStringInterface(t.Args)
	p := bluemonday.UGCPolicy()

	//Check deoes Username, and email provided
	if args["name"] == nil || args["mail"] == nil {
		return ErrorReturn(t, 404, "000012", "Please provide name and mail")
	}

	unmae := p.Sanitize(fmt.Sprintf("%v", args["name"]))
	umail := p.Sanitize(fmt.Sprintf("%v", args["mail"]))

	db, err := ConnectDBv2()
	if err != nil {
		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			SetErrorLog(err.Error())
		}

		return ErrorReturn(t, 500, "000027", err.Error())
	}

	maxusers := viper.GetInt("settings.user_creation")
	domain := viper.GetString("server.domain")

	curdata := Users{}
	var count int64
	db.Conn.Debug().Model(curdata).Count(&count)

	if count >= int64(maxusers) {
		return ErrorReturn(t, 403, "000027", "You are not alowed to create more users")
	}

	newdata := Users{}
	newdatainfo := UsersInfo{}

	//Check does data is unique
	var countb int64
	db.Conn.Debug().Where("name = ? OR mail = ?", unmae, umail).Model(curdata).Count(&countb)
	if countb != 0 {
		return ErrorReturn(t, 403, "000027", "User with such name and email is already exist")
	}

	JsonArgs, err := json.Marshal(args)
	if err != nil {
		return ErrorReturn(t, 500, "000028", err.Error())
	}

	err = json.Unmarshal(JsonArgs, &newdata)
	if err != nil {
		return ErrorReturn(t, 500, "000028", err.Error())
	}

	err = json.Unmarshal(JsonArgs, &newdatainfo)
	if err != nil {
		return ErrorReturn(t, 500, "000028", err.Error())
	}

	userid := Hashgen(8)
	newdata.UID = userid
	newdatainfo.UID = userid
	password := Hashgen(14)

	//Check does password provided
	if args["pass"] != nil {

		password = fmt.Sprintf("%v", args["pass"])
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {

		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			SetErrorLog("dbstructure.go: " + err.Error())
		}
	}

	newdata.Pass = string(hashedPassword)
	newdata.Mailsent = int(time.Now().Unix())
	newdata.Created = int(time.Now().Unix())
	newdata.Mailconfirmed = int(time.Now().Unix())
	newdata.Status = true
	newdata.Completed = true

	err = db.Conn.Create(&newdata).Error
	if err != nil {
		return ErrorReturn(t, 400, "000005", err.Error())
	}

	err = db.Conn.Create(&newdatainfo).Error
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
	msga = fmt.Sprintf("login: %s", unmae)
	message = append(message, msga)
	msga = fmt.Sprintf("password: <code>%s</code>", password)
	message = append(message, msga)
	msga = "Do not forget to change your password after first login!"
	message = append(message, msga)

	go SendNotification(t, title, message, "invitation", userid)

	ans["uid"] = userid
	return Interfacetoresponse(t, ans)

}
