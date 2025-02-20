package middleware

import (
	"crypto/sha512"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

var cookieSalt string

func SetCookieSalt(salt string) {
	cookieSalt = salt
}

func GinSetCookieID(c *gin.Context, id int) {
	id_str := strconv.Itoa(id)
	c.SetCookie("user_id", id_str, 3600, "/", "", false, true)

	sign := cookieSign([]string{id_str})
	c.SetCookie("sign", sign, 3600, "/", "", false, true)
}

func GinCheckRequesCookie(c *gin.Context) bool {
	id, id_err := c.Cookie("user_id")
	sign, sign_err := c.Cookie("sign")
	if id_err == nil && sign_err == nil {
		return isCookieValid(sign, []string{id})
	}
	return false
}

func isCookieValid(sign string, cookies []string) bool {
	calculateHash := cookieSign(cookies)
	return calculateHash == sign
}

func cookieSign(cookies []string) (sign string) {
	result_string := ""
	for _, s := range cookies {
		result_string += s
	}
	result_string += cookieSalt

	b := sha512.Sum512([]byte(result_string))
	sign = fmt.Sprintf("%x", b)

	return
}
