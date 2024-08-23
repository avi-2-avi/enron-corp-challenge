package utils

import "encoding/base64"

func GetAuthHeader() string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte("admin:Pass123!!!"))
}
