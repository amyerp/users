//////////////////////////////////////////////////////////////////////////////////
// Copyright 2021 Alexey Yanchenko <mail@yanchenko.me>                          //
//                                                                              //
// This file is part of the ERP library.                                        //
//                                                                              //
//  Unauthorized copying of this file, via any media is strictly prohibited     //
//  Proprietary and confidential                                                //
//////////////////////////////////////////////////////////////////////////////////

package functions

import "time"

func FormatDate(date time.Time, dateformat string) string {

	if date.IsZero() {
		return ""
	}

	formateddate := date.Format(dateformat)
	return formateddate

}
