package main

import (
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	url2 "net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	headers := make(map[string][]string)
	forms := make(map[string][]string)
	dataRaw := ""
	//filePath
	dataBinary := ""
	method := flag.String("X", "GET", "request method ,default value is get")
	requestUrl := flag.String("h", "", "request url")
	timeout := flag.Int("t", 5000, "request time out,unit is ms ,default 5000 ms")
	responseHeader := flag.Bool("i", true, "println response header,default true")
	flag.Parse()

	args := flag.Args()
	isMultiForm := false
	if len(args) > 0 {
		for i, arg := range args {
			switch arg {
			case "--header":
				fallthrough
			case "-H":
				header := args[i+1]
				if strings.Contains(header, ":") {
					kv := strings.Split(header, ":")
					headers[kv[0]] = []string{kv[1]}
				}
				continue
			case "--form":
				form := args[i+1]
				if strings.Contains(form, "=") {
					kv := strings.Split(form, "=")
					formValue := kv[1]
					if strings.HasPrefix(formValue, "\"") {
						formValue = formValue[1:]
					}
					if strings.HasSuffix(formValue, "\"") {
						formValueLen := len(formValue)
						formValue = formValue[:formValueLen-1]
					}
					if !isMultiForm {
						isMultiForm = strings.HasPrefix(formValue, "@")
					}
					forms[kv[0]] = []string{formValue}
				}
				continue
			case "--data-raw":
				fallthrough
			case "-d":
				dataRaw = args[i+1]
			case "--data-binary":
				dataBinary = arg
				log.Panic("暂不支持", dataBinary)
			default:
				if *requestUrl == "" {
					requestUrl = &arg
				}
			}
		}
		requestUrl = &args[0]
	}

	if *requestUrl == "" {
		log.Panic("url不能为空")
	}
	timeoutDuration, _ := time.ParseDuration(strconv.Itoa(*timeout) + "ms")
	client := &http.Client{
		Timeout: timeoutDuration,
	}
	if url, err := url2.Parse(*requestUrl); err != nil {
		log.Panicf("url[%s] unavaliable", url)
	} else {
		request, _ := http.NewRequest(*method, *requestUrl, strings.NewReader(dataRaw))
		request.Form = forms
		if dataRaw != "" {
			headers["Content-Type"] = []string{"application/json"}
		}
		if len(forms) > 0 {
			headers["Content-Type"] = []string{"application/x-www-form-urlencoded"}
		}

		request.Header = headers

		if isMultiForm {
			multipartRequest(request, forms)
		}
		if resp, err := client.Do(request); err != nil {
			log.Panicf("url[%s] error,%v", url, err)
		} else {
			//打印请求头
			if *responseHeader {
				for k, v := range resp.Header {
					log.Println(k, ":", v)
				}
			}
			data, _ := ioutil.ReadAll(resp.Body)
			log.Printf("%s", data)
		}
	}
}

func multipartRequest(request *http.Request, forms map[string][]string) {
	body := &bytes.Buffer{}
	fileWriter := multipart.NewWriter(body)
	defer fileWriter.Close()
	for key, val := range forms {
		if strings.HasPrefix(val[0], "@") {
			uploadFilePath := val[0][1:]
			if file, err := os.Open(uploadFilePath); err != nil {
				log.Println("无法打开的文件:", uploadFilePath, err.Error())
			} else if part, err := fileWriter.CreateFormFile(key, uploadFilePath); err != nil {
				log.Println("无法识别的文件:", uploadFilePath, err.Error())
			} else if _, err = io.Copy(part, file); err != nil {
				log.Println("无法读取的文件:", uploadFilePath, err.Error())
			} else {
				log.Printf("文件[%s]被成功读取\n", uploadFilePath)
			}
		}
		_ = fileWriter.WriteField(key, val[0])
	}
	request.Header.Set("Content-Type", fileWriter.FormDataContentType())
}
