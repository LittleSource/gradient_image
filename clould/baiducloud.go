package clould

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var Header = map[string]string{"Content-Type": "application/x-www-form-urlencoded"}

type BaiduClould struct {
	clientId     string
	clientSecret string
	AccessToken  string
}

func NewBaiduClould(clientId, clientSecret string) *BaiduClould {
	accessToken, err := getAccessToken(clientId, clientSecret)
	if err != nil {
		panic(err)
	}
	return &BaiduClould{
		clientId:     clientId,
		clientSecret: clientSecret,
		AccessToken:  accessToken,
	}
}

func (b *BaiduClould) PortraitSegmentation(reader io.Reader, outputPath string) error {
	var host = "https://aip.baidubce.com/rest/2.0/image-classify/v1/body_seg"
	uri, err := url.Parse(host)
	if err != nil {
		return err
	}
	query := uri.Query()
	query.Set("access_token", b.AccessToken)
	uri.RawQuery = query.Encode()
	filebytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	image := base64.StdEncoding.EncodeToString(filebytes)
	sendBody := http.Request{}
	sendBody.ParseForm()
	sendBody.Form.Add("image", image)
	result, err := httpPost(uri.String(), sendBody.Form.Encode(), Header)
	if err != nil {
		return err
	}
	resp := struct {
		Labelmap   string `json:"labelmap"`
		Scoremap   string `json:"scoremap"`
		Foreground string `json:"foreground"`
		PersonNum  int32  `json:"person_num"`
	}{}
	err = json.Unmarshal(result, &resp)
	if err != nil {
		return err
	}
	data, err := base64.StdEncoding.DecodeString(resp.Scoremap)
	if err != nil {
		return err
	}
	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	out.Write(data)
	out.Close()
	return nil
}

func getAccessToken(clientId, clientSecret string) (string, error) {
	var host = "https://aip.baidubce.com/oauth/2.0/token"
	sendBody := http.Request{}
	sendBody.ParseForm()
	sendBody.Form.Add("grant_type", "client_credentials")
	sendBody.Form.Add("client_id", clientId)
	sendBody.Form.Add("client_secret", clientSecret)
	result, err := httpPost(host, sendBody.Form.Encode(), Header)
	if err != nil {
		return "", err
	}
	at := struct {
		AccessToken string `json:"access_token"`
	}{}
	err = json.Unmarshal(result, &at)
	if err != nil {
		return "", err
	}
	return at.AccessToken, nil
}

func httpPost(uri string, sendData string, headers ...map[string]string) ([]byte, error) {
	request, err := http.NewRequest("POST", uri, strings.NewReader(sendData))
	for _, header := range headers {
		for k, v := range header {
			request.Header.Set(k, v)
		}
	}
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return result, nil
}
