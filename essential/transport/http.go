package transport

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Cookies resp cookies
var Cookies []*http.Cookie

// ReqAPI func
func ReqAPI(api string, params map[string]string, payload interface{}) (respmsg Response, err error) {

	requrl := moduleConfig.Server + api
	urlObj, err := url.Parse(requrl)
	if err != nil {
		return respmsg, err
	}
	querys := urlObj.Query()

	for key, value := range params {
		querys.Set(key, value)
	}

	urlObj.RawQuery = querys.Encode()
	connurl := urlObj.String()

	reqbody, err := json.Marshal(payload)
	if err != nil {
		return respmsg, err
	}
	// 通过 http 请求
	req, _ := http.NewRequest("GET", connurl, bytes.NewReader(reqbody))
	for _, cookie := range Cookies {
		req.AddCookie(cookie)
	}

	resp, err := DefaultClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return respmsg, err
	}

	if resp.StatusCode != 200 {
		errmsg := fmt.Sprintf("api not available, status code: %d", resp.StatusCode)
		return respmsg, errors.New(errmsg)
	}

	Cookies = resp.Cookies()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return respmsg, err
	}
	err = json.Unmarshal(body, &respmsg)
	if err != nil {
		return respmsg, err
	}

	return respmsg, nil
}

// PostAPI func
func PostAPI(api string, params map[string]string, payload interface{}) (respmsg Response, err error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return respmsg, err
	}
	requrl := moduleConfig.Server + api
	// post api
	querys := url.Values{}

	urlObj, err := url.Parse(requrl)
	if err != nil {
		return respmsg, err
	}

	for key, value := range params {
		querys.Set(key, value)
	}

	urlObj.RawQuery = querys.Encode()
	connurl := urlObj.String()

	request, err := http.NewRequest("POST", connurl, bytes.NewReader(data))
	if err != nil {
		return respmsg, err
	}
	for _, cookie := range Cookies {
		request.AddCookie(cookie)
	}

	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	resp, err := DefaultClient.Do(request)
	// 这里需要处理出错机制，比如上传失败，需要重新上传等，这里后续完善
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return respmsg, err
	}

	if resp.StatusCode != 200 {
		errmsg := fmt.Sprintf("api not available, status code: %d", resp.StatusCode)
		return respmsg, errors.New(errmsg)
	}
	Cookies = resp.Cookies()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return respmsg, err
	}

	err = json.Unmarshal(body, &respmsg)
	if err != nil {
		return respmsg, err
	}

	return respmsg, nil
}

// DelAPI func
func DelAPI(api string, params map[string]string) (respmsg Response, err error) {

	requrl := moduleConfig.Server + api
	// post api
	querys := url.Values{}

	urlObj, err := url.Parse(requrl)
	if err != nil {
		return respmsg, err
	}

	for key, value := range params {
		querys.Set(key, value)
	}

	urlObj.RawQuery = querys.Encode()
	connurl := urlObj.String()

	request, err := http.NewRequest("DELETE", connurl, nil)
	if err != nil {
		return respmsg, err
	}

	for _, cookie := range Cookies {
		request.AddCookie(cookie)
	}

	resp, err := DefaultClient.Do(request)
	// 这里需要处理出错机制，比如上传失败，需要重新上传等，这里后续完善
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return respmsg, err
	}

	if resp.StatusCode != 200 {
		errmsg := fmt.Sprintf("api not available, status code: %d", resp.StatusCode)
		return respmsg, errors.New(errmsg)
	}
	Cookies = resp.Cookies()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return respmsg, err
	}

	err = json.Unmarshal(body, &respmsg)
	if err != nil {
		return respmsg, err
	}

	return respmsg, nil
}
