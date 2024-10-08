//////////////////////////////////////////////////////////////////////////////////
// Copyright 2021 Alexey Yanchenko <mail@yanchenko.me>                          //
//                                                                              //
// This file is part of the ERP library.                                        //
//                                                                              //
//  Unauthorized copying of this file, via any media is strictly prohibited     //
//  Proprietary and confidential                                                //
//////////////////////////////////////////////////////////////////////////////////

package model

import (
	"gorm.io/gorm"
)

type APITokens struct {
	gorm.Model
	TokenId    string `gorm:"column:tokenid;type:varchar(60);UNIQUE;NOT NULL;"  json:"tokenid"`
	Token      string `gorm:"column:token;type:varchar(254);UNIQUE;NOT NULL;"  json:"token"`
	TokenName  string `gorm:"column:tokenname;type:varchar(60);DEFAULT '';" json:"tokenname"`
	UID        string `gorm:"column:uid;type:varchar(60);NOT NULL;" json:"uid"`
	Created    int    `gorm:"column:created;type:int;DEFAULT '0'" json:"created"`
	Expiration int    `gorm:"column:expiration;type:int;DEFAULT '0'" json:"expiration"`  //if 0 - no expiration time
	Status     bool   `gorm:"column:status;type:bool;DEFAULT 'true'" json:"status"`      // if true - active, if false - deactivated
	IsAdmin    bool   `gorm:"column:is_admin;type:bool;DEFAULT 'false'" json:"is_admin"` //only if generated by admin
	Readonly   bool   `gorm:"column:readonly;type:bool;DEFAULT 'false'" json:"readonly"`
	Comment    string `gorm:"column:comment;type:varchar(60);DEFAULT '';" json:"comment"`
}

type Users struct {
	gorm.Model
	UID           string `gorm:"column:uid;type:varchar(60);UNIQUE;NOT NULL;" json:"uid"` //userID
	Pass          string `gorm:"column:pass;type:varchar(128);NOT NULL;DEFAULT ''" json:"pass,omitempty"`
	Name          string `gorm:"column:name;type:varchar(60);NOT NULL;DEFAULT '';UNIQUE" json:"name,omitempty"`
	Mail          string `gorm:"column:mail;type:varchar(254);DEFAULT '';UNIQUE"  json:"mail,omitempty"`
	Mailsent      int    `gorm:"column:mailsent;type:int;DEFAULT '0'" json:"mailsent,omitempty"`
	Mailconfirmed int    `gorm:"column:mailconfirmed;:int;DEFAULT '0'" json:"mailconfirmed,omitempty"`
	Created       int    `gorm:"column:created;type:int;DEFAULT '0'" json:"created,omitempty"`
	Access        int    `gorm:"column:access;type:int;DEFAULT '0'" json:"access,omitempty"`
	Login         int    `gorm:"column:login;type:int;DEFAULT '0'" json:"login"`
	IP            string `gorm:"column:ip;type:varchar(128); DEFAULT ''" json:"ip,omitempty"`
	Status        bool   `gorm:"column:status;type:bool;DEFAULT 'false'" json:"status,omitempty"`
	Completed     bool   `gorm:"column:completed;type:bool;DEFAULT 'false'" json:"completed,omitempty"`
	IsAdmin       bool   `gorm:"column:is_admin;type:bool;DEFAULT 'false'" json:"isadmin"`
	Readonly      bool   `gorm:"column:readonly;type:bool;DEFAULT 'false'" json:"readonly,omitempty"`
	TFA           bool   `gorm:"column:tfa;type:bool;DEFAULT false;" json:"tfa"`
	TFAType       string `gorm:"column:tfatype;type:varchar(60);DEFAULT '';" json:"tfatype,omitempty"`
}

type UsersInfo struct {
	gorm.Model
	UID         string `gorm:"column:uid;type:varchar(60);UNIQUE;NOT NULL;" json:"user_id,omitempty"` //userID
	PersonID    string `gorm:"column:personid;type:varchar(60);DEFAULT '';" json:"personid,omitempty"`
	Name        string `gorm:"column:name;type:varchar(254);DEFAULT '';" json:"first_name,omitempty"`
	MName       string `gorm:"column:mname;type:varchar(254);DEFAULT '';" json:"middle_name,omitempty"`
	Surname     string `gorm:"column:surname;type:varchar(254);DEFAULT '';" json:"surname,omitempty"`
	AvatarID    string `gorm:"column:avatarid;type:varchar(254);DEFAULT '';" json:"avatarid,omitempty"` //fileid in files table
	BirthDate   string `gorm:"column:birthdate;type:varchar(60);DEFAULT '';" json:"birthdate,omitempty"`
	PhoneNumber string `gorm:"column:phonenumber;type:varchar(254);DEFAULT '';" json:"phonenumber,omitempty"`
	Role        string `gorm:"column:role;type:varchar(254);DEFAULT '';" json:"role,omitempty"`
}

type UserSettings struct {
	gorm.Model
	UID        string `gorm:"column:uid;type:varchar(60);UNIQUE;NOT NULL;" json:"uid"` //userID
	DateFormat string `gorm:"column:dateformat;type:varchar(60);DEFAULT '2006-01-02';" json:"dateformat"`
	NightMode  bool   `gorm:"column:night_mode;type:bool;DEFAULT 'false'" json:"night_mode"`
}

/*
type Users struct {
	gorm.Model
	UID      string `gorm:"column:uid;type:varchar(60);UNIQUE;NOT NULL;" json:"uid"` //userID
	Name     string `gorm:"column:name;type:varchar(60);NOT NULL;DEFAULT '';UNIQUE" json:"name"`
	Mail     string `gorm:"column:mail;type:varchar(254);DEFAULT '';UNIQUE"  json:"mail"`
	Created  int    `gorm:"column:created;type:int;DEFAULT '0'" json:"created"`
	IsAdmin  bool   `gorm:"column:is_admin;type:bool;DEFAULT 'false'" json:"isadmin"`
	Readonly bool   `gorm:"column:readonly;type:bool;DEFAULT 'false'" json:"readonly"`
	PersonID string `gorm:"column:personid;type:varchar(60);DEFAULT '';" json:"personid"`
	Role     string `gorm:"column:role;type:varchar(60);DEFAULT '';" json:"role"`
	TFA      bool   `gorm:"column:tfa;type:bool;DEFAULT false;" json:"tfa"`
	TFAType  string `gorm:"column:tfatype;type:varchar(60);DEFAULT '';" json:"tfatype"`
}

type UserSettings struct {
	gorm.Model
	UID        string `gorm:"column:uid;type:varchar(60);UNIQUE;NOT NULL;" json:"uid"` //userID
	DateFormat string `gorm:"column:dateformat;type:varchar(60);DEFAULT '2006-01-02';" json:"dateformat"`
	NightMode  string `gorm:"column:night_mode;type:varchar(60);DEFAULT '2006-01-02';" json:"night_mode"`
}
*/

type UserResponse struct {
	*Users
	*UsersInfo
}
