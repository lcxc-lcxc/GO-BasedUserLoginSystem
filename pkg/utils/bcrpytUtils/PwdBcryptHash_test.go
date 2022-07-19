/**
 @author: 15973
 @date: 2022/07/19
 @note:
**/
package bcrpytUtils

import (
	"fmt"
	"testing"
)

func TestPwdHash(t *testing.T) {

	verify := PwdVerify("$2a$10$L/JSt4kpkHjBWKjSazFX.eeDyzbk5qp37ZcX3SHJUW841yGQt5BC2", "1234567")
	fmt.Println(verify)
	//time.Sleep(time.Second)
	//c, _ := PwdHash("1234567")
	//time.Sleep(time.Second)
	//d, _ := PwdHash("1234567")
	//t.Log(PwdVerify(a, "1234567"))
	//t.Log(PwdVerify(b, "1234567"))
	//t.Log(PwdVerify(c, "1234567"))
	//t.Log(PwdVerify(d, "1234567"))

}
