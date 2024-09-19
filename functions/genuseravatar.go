package functions

import (
	"image/png"
	"os"
	"path/filepath"
	"time"
	. "user/model"

	"github.com/getsentry/sentry-go"
	. "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/spf13/viper"
)

func GenUserAvatar(userid string, t *pb.Request) (link string) {

	fileid := Hashgen(12)
	extension := ".png"

	var pwd string = viper.GetString("server.filedir")

	pathfiles := filepath.Join(pwd, "users")
	//Create dir output using above code
	if _, err := os.Stat(pathfiles); os.IsNotExist(err) {
		os.Mkdir(pathfiles, 0755)
	}

	pathfiles = filepath.Join(pathfiles, userid)
	//Create dir output using above code
	if _, err := os.Stat(pathfiles); os.IsNotExist(err) {
		os.Mkdir(pathfiles, 0755)
	}

	pathfiles = filepath.Join(pathfiles, "avatar")
	//Create dir output using above code
	if _, err := os.Stat(pathfiles); os.IsNotExist(err) {
		os.Mkdir(pathfiles, 0755)
	}

	filelink := pathfiles + "/" + fileid + extension

	out, err := os.Create(filelink)
	if err != nil {
		SetLog("gravatar " + err.Error())
	}

	timestamp := time.Now().String()
	avahash := userid + timestamp
	png.Encode(out, CreateGravatar([]byte(avahash), 540, 60))

	//Add Link to product table
	db, err := ConnectDBv2()
	if err != nil {
		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			SetErrorLog(err.Error())
		}

	}

	user := &UsersInfo{}
	user.AvatarID = filelink
	db.Conn.Model(&user).Where("uid = ?", userid).Updates(&user)

	return filelink

}
