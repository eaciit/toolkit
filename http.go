package toolkit

/*
Http related
*/
import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func HttpCall(url string, callType string, datas []M,
	useAuth bool, userName string, password string) (*http.Response, error) {
	var err error

	//-- preparing cookie jar and http client
	/*
		jar, err := cookiejar.New(nil)
		client := &http.Client{
			Jar: jar,
		}
	*/
	client := &http.Client{}

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
	/*
		if dbSessionId != "" {
			expire := time.Time{}
			cookieSession := http.Cookie{"OSESSIONID", dbSessionId, "/", url,
				expire, expire.Format(time.UnixDate), 86400, true, true,
				"OSESSIONID=" + dbSessionId, []string{"OSESSIONID=" + dbSessionId}}
			req.AddCookie(&cookieSession)
		}
	*/
	resp, err = client.Do(req)
	return resp, err
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
