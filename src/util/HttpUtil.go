package util

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

func GetRequest(url string) (string, error) {
	client := &http.Client{}
	reqest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	response, err := client.Do(reqest)
	if err != nil {
		return "", err
	}
	if response.StatusCode != 200 {
		return "", errors.New("response code is " + response.Status)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

var client = &http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) == 200 {
			return errors.New("stopped after 200 redirects")
		}
		return nil
	},
}

func GetRequestForReader(url string) (io.Reader, error) {
	reqest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	reqest.Header.Add("Cookie", "bdshare_firstime=1534401393293; __cfduid=daf43bbb0f281acf06c165d1dd53d96081540620552; Hm_lvt_dca7dc99d8ac55393ef7fbc057d85ffb=1540620586,1540620630; PHPSESSID=51tls5gckf7imtpbc2l23fpt82; Hm_lvt_40108d7e4cc326e04eecdd70da888247=1550040867; damon_token_2019=679885681; Murl=-m0XB7fGVrb8WLe%3DkTv8V3XF%2F7v2NypDU0yRAYh0N7tDWYVKkq2_A2q8xr_BVRvS6RPdFZP0N4s0B; Hm_lpvt_40108d7e4cc326e04eecdd70da888247=1550045323")
	response, err := client.Do(reqest)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New("response code is " + response.Status)
	}
	return response.Body, nil
}
