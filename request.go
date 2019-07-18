package agollo

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type CallBack struct {
	SuccessCallBack   func([]byte) (interface{}, error)
	NotModifyCallBack func() error
}

type ConnectConfig struct {
	//设置到http.client中timeout字段
	Timeout time.Duration
	//连接接口的uri
	Uri string
}

var client = &http.Client{
	Timeout: 2 * time.Minute,
	Transport: &http.Transport{
		ResponseHeaderTimeout: 2 * time.Minute,
	},
}

func request(requestUrl string, connectionConfig *ConnectConfig, callBack *CallBack) (value interface{}, err error) {
	var responseBody []byte
	var res *http.Response

	res, err = client.Get(requestUrl)
	if res == nil || err != nil {
		logger.Println("Connect Apollo Server Fail,Error:", err)
		return
	}

	defer res.Body.Close()

	//not modified break
	switch res.StatusCode {
	case http.StatusOK:
		responseBody, err = ioutil.ReadAll(res.Body)
		if err != nil {
			logger.Println("Connect Apollo Server Fail,Error:", err)
			return
		}

		logger.Println("[DEBUG]: ", string(responseBody))

		if callBack != nil && callBack.SuccessCallBack != nil {
			return callBack.SuccessCallBack(responseBody)
		} else {
			return nil, nil
		}
	case http.StatusNotModified:
		logger.Println("Config Not Modified:", err)
		if callBack != nil && callBack.NotModifyCallBack != nil {
			return nil, callBack.NotModifyCallBack()
		} else {
			return nil, nil
		}
	default:
		logger.Println("Connect Apollo Server Fail,Error:", err)
		if res != nil {
			logger.Println("Connect Apollo Server Fail,StatusCode:", res.StatusCode)
		}
		err = errors.New("connect Apollo Server Fail")
	}

	return
}

func requestRecovery(appConfig *AppConfig,
	connectConfig *ConnectConfig,
	callBack *CallBack) (interface{}, error) {
	format := "%s%s"
	var err error
	var response interface{}

	for {
		host := appConfig.selectHost()
		if host == "" {
			return nil, err
		}

		requestUrl := fmt.Sprintf(format, host, connectConfig.Uri)
		response, err = request(requestUrl, connectConfig, callBack)
		if err == nil {
			return response, err
		}

		setDownNode(host)
	}

}
