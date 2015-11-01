package toolkit

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func httpCall(url string, callType string, datas []T,
	useAuth bool, userName string, password string) (*http.Response, error) {
	//-- preparing cookie jar and http client
	jar, err := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	//-- GET
	var resp *http.Response
	var req *http.Request
	if callType == "GET" {
		resp, err = client.Get(url)
		req, err = http.NewRequest(callType, url, nil)
	}
	if useAuth == true {
		req.SetBasicAuth(userName, password)
	}
	if dbSessionId != "" {
		expire := time.Time{}
		cookieSession := http.Cookie{"OSESSIONID", dbSessionId, "/", url,
			expire, expire.Format(time.UnixDate), 86400, true, true,
			"OSESSIONID=" + dbSessionId, []string{"OSESSIONID=" + dbSessionId}}
		req.AddCookie(&cookieSession)
	}
	resp, err = client.Do(req)
	return resp, err
}

func HttpContent(r *http.Response) []byte {
	defer r.Body.Close()
	bytes, _ := ioutil.ReadAll(r.Body)
	return bytes
}

func HttpContentM(r *http.Response) M {
	bytes := getContent(r)
	obj := M{}
	_ = json.Unmarshal(bytes, &obj)
	return obj
}

func HttpContentString(r *http.Response) Obj {
	bytes := getContent(r)
	obj := string(bytes)
	return obj
}
