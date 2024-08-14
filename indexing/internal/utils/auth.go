package utils

import "encoding/base64"

func GetAuthHeader() string {
	auth := base64.StdEncoding.EncodeToString([]byte("admin:Pass123!!!"))
	return "Basic " + auth
}
