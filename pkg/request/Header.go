package request

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/VeronicaAlexia/pineapple-backups/config"
	"net/http"
	"os"
)

func Base64Bytes(UserName, Password string) string {
	var encoded bytes.Buffer
	authentication := []byte(UserName + "&" + Password)
	if _, err := base64.NewEncoder(base64.StdEncoding, &encoded).Write(authentication); err == nil {
		return string(encoded.Bytes())
	} else {
		fmt.Println("encoder.Write:", err)
	}
	return ""

}

func SET_THE_HEADERS(req *http.Request) {
	HeaderCollection := make(map[string]string)
	HeaderCollection["Content-Type"] = "application/json"
	switch config.Vars.AppType {
	case "sfacg":
		HeaderCollection["sf-minip-info"] = "minip_novel/1.0.70(android;11)/wxmp"
		HeaderCollection["Cookie"] = config.Apps.Sfacg.Cookie
		HeaderCollection["Authorization"] = Base64Bytes(config.Apps.Sfacg.UserName, config.Apps.Sfacg.Password)
	case "cat":
		HeaderCollection["User-Agent"] = "Android  com.kuangxiangciweimao.novel  2.9.291, Google, Pixel5"
		HeaderCollection["Cookie"] = "Account:" + config.Apps.Cat.Params.Account + ";" + config.Apps.Cat.Params.LoginToken
		HeaderCollection["Authorization"] = Base64Bytes(config.Apps.Cat.Params.Account, config.Apps.Cat.Params.LoginToken)

	default:
		fmt.Println(config.Vars.AppType, "AppType is invalid, please check config file")
		os.Exit(1)
	}
	for HeaderKey, HeaderValue := range HeaderCollection {
		req.Header.Set(HeaderKey, HeaderValue)

	}
}