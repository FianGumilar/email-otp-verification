package utils

import (
	"encoding/json"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/FianGumilar/email-otp-verification/constans"
	"github.com/FianGumilar/email-otp-verification/models"
	"github.com/labstack/echo"
)

func QueryFill(query string) (new string) {
	query = strings.ReplaceAll(query, " ", "")
	split := strings.Split(query, ",")
	for range split {
		new += "?,"
	}

	return strings.TrimSuffix(new, ",")
}

func ReplaceSQL(old, searchPattern string) string {
	tmpCount := strings.Count(old, searchPattern)
	for m := 1; m <= tmpCount; m++ {
		old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(m), 1)
	}
	return old
}

func BindValidateStruct(ctx echo.Context, i interface{}, function string) error {
	if err := ctx.Bind(i); err != nil {
		return err
	}
	bytes, _ := json.Marshal(i)
	log.Println("Incoming Request on", function, "=>", string(bytes))
	if err := ctx.Validate(i); err != nil {
		return err
	}

	return nil
}

func ResponseJSON(success bool, code string, msg string, result interface{}) models.Response {
	tm := time.Now()
	response := models.Response{
		Success:          success,
		StatusCode:       code,
		Result:           result,
		Message:          msg,
		ResponseDatetime: tm,
	}

	return response
}

func TimeStampNow() string {
	return time.Now().Format(constans.LAYOUT_TIMESTAMP)
}

func GenerateRandomString(n int) string {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	charsets := []rune("abcdefghijklmnopqrzABCDEFGHIJKLM234567890")
	letters := make([]rune, n)
	for i := range letters {
		letters[i] = charsets[r.Intn(len(charsets))]
	}
	return string(letters)
}

func GenerateRandomNumber(n int) string {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	charsets := []rune("0123456789")
	letters := make([]rune, n)
	for i := range letters {
		letters[i] = charsets[r.Intn(len(charsets))]
	}
	return string(letters)
}
