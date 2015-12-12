package toolkit

/*
Http related
*/
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

func HttpCall(url string, callType string,
	datas []byte,
	config M) (*http.Response, error) {
	var err error
	if config == nil {
		config = M{}
	}

	var req *http.Request

	//-- GET
	if datas == nil || len(datas) == 0 {
		req, err = http.NewRequest(callType, url, nil)
	} else {
		//fmt.Printf("Datas: %v\n%s \n", datas, string(datas))
		rdr := bytes.NewBuffer(datas)
		req, err = http.NewRequest(callType, url, rdr)
	}
	if err != nil {
		return nil, err
	}

	if config.Has("auth") {
		authType := strings.ToLower(config.Get("auth", "").(string))
		if authType == "basic" {
			username := config.Get("user", "").(string)
			pass := config.Get("password", "").(string)
			req.SetBasicAuth(username, pass)
		}
	}

	return httpcall(req, config)
}

func httpcall(req *http.Request, config M) (*http.Response, error) {
	var client *http.Client

	//-- handling cookie
	if config.Has("cookie") == false {
		client = new(http.Client)
	} else {
		//-- preparing cookie jar and http client
		jar, err := cookiejar.New(nil)
		if err != nil {
			return nil, fmt.Errorf("Unable to initialize cookie jar: %s", err.Error())
		}
		client = &http.Client{
			Jar: jar,
		}
	}

	//--- handling header
	if headers, hasHeaders := config["headers"]; hasHeaders {
		mHeaders := headers.(M)
		for k, v := range mHeaders {
			req.Header.Add(k, v.(string))
		}
	}

	var resp *http.Response
	var errCall error
	resp, errCall = client.Do(req)
	if expectedStatus := config.Get("expectedstatus", 0).(int); expectedStatus != 0 && resp.StatusCode != expectedStatus {
		return nil, fmt.Errorf("Code error: " + resp.Status)
	}
	return resp, errCall
}

func HttpContent(r *http.Response) []byte {
	defer r.Body.Close()
	bytes, _ := ioutil.ReadAll(r.Body)
	return bytes
}

func HttpContentM(r *http.Response) M {
	bytes := HttpContent(r)
	obj := M{}
	_ = json.Unmarshal(bytes, &obj)
	return obj
}

func HttpContentString(r *http.Response) string {
	bytes := HttpContent(r)
	return string(bytes)
}
