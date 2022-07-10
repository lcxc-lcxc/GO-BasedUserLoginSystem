/**
 @author: 15973
 @date: 2022/07/10
 @note:
**/
package test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"v0.0.0/config"
	"v0.0.0/internel/api/user"
)

func TestUserRegister(t *testing.T) {

	resp, err := http.Post("http://"+config.HttpAddress+"/api/user/register",
		"application/x-www-form-urlencoded",
		strings.NewReader("username=123456&password=123456&nickname=123456"))
	if err != nil {
		t.Errorf("TestUserRegister: request call failed : %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("TestUserRegister: response body open failed : %v", err)
	}
	registerResp := &user.RegisterResponse{}
	err = json.Unmarshal(body, registerResp)
	if err != nil {
		t.Errorf("TestUserRegister: response body unmarshal failed : %v", err)
	}
	t.Logf("%#v", registerResp)

}
