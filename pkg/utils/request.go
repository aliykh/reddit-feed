package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"

	JSON = "json"
	FORM = "form"
)

var (
	ErrMIMENotSupported = errors.New("mime is not supported")
)

type Response struct {
	StatusCode int
	Body       []byte
}

func InvokeHandler(req *http.Request, router *gin.Engine) (resp *Response, err error) {

	// initialize response record
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// extract the response from the response record
	result := w.Result()
	defer result.Body.Close()

	// extract response body
	bodyByte, err := ioutil.ReadAll(result.Body)

	resp = &Response{
		StatusCode: result.StatusCode,
		Body:       bodyByte,
	}

	return
}

func MakeRequest(method, mime, api string, param interface{}) (request *http.Request, err error) {
	method = strings.ToUpper(method)
	mime = strings.ToLower(mime)

	switch mime {
	case JSON:
		var (
			contentBuffer *bytes.Buffer
			jsonBytes     []byte
		)
		jsonBytes, err = json.Marshal(param)
		if err != nil {
			return
		}
		contentBuffer = bytes.NewBuffer(jsonBytes)
		request, err = http.NewRequest(method, api, contentBuffer)
		if err != nil {
			return
		}
		request.Header.Set("Content-Type", "application/json;charset=utf-8")
	case FORM:
		queryStr := MakeQueryStrFrom(param)
		var buffer io.Reader

		if (method == DELETE || method == GET) && queryStr != "" {
			api += "?" + queryStr
		} else {
			buffer = bytes.NewReader([]byte(queryStr))
		}

		request, err = http.NewRequest(string(method), api, buffer)
		if err != nil {
			return
		}
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	default:
		err = ErrMIMENotSupported
		return
	}
	return
}
func MakeQueryStrFrom(params interface{}) (result string) {
	if params == nil {
		return
	}
	value := reflect.ValueOf(params)

	switch value.Kind() {
	case reflect.Struct:
		var formName string
		for i := 0; i < value.NumField(); i++ {
			if formName = value.Type().Field(i).Tag.Get("form"); formName == "" {
				// don't tag the form name, use camel name
				formName = GetCamelNameFrom(value.Type().Field(i).Name)
			}
			result += "&" + formName + "=" + fmt.Sprintf("%v", value.Field(i).Interface())
		}
	case reflect.Map:
		for _, key := range value.MapKeys() {
			result += "&" + fmt.Sprintf("%v", key.Interface()) + "=" + fmt.Sprintf("%v", value.MapIndex(key).Interface())
		}
	default:
		return
	}

	if result != "" {
		result = result[1:]
	}
	return
}
func GetCamelNameFrom(name string) string {
	result := ""
	i := 0
	j := 0
	r := []rune(name)
	for m, v := range r {
		// if the char is the capital
		if v >= 'A' && v < 'a' {
			// if the prior is the lower-case || if the prior is the capital and the latter is the lower-case
			if (m != 0 && r[m-1] >= 'a') || ((m != 0 && r[m-1] >= 'A' && r[m-1] < 'a') && (m != len(r)-1 && r[m+1] >= 'a')) {
				i = j
				j = m
				result += name[i:j] + "_"
			}
		}
	}

	result += name[j:]
	return strings.ToLower(result)
}
