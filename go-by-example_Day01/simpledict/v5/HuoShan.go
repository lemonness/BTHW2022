package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type HuoShanDictRequest struct {
	Text     string `json:"text"`
	Language string `json:"language"`
}

type HuoShanDictResponse struct {
	Words []struct {
		Source  int    `json:"source"`
		Text    string `json:"text"`
		PosList []struct {
			Type      int `json:"type"`
			Phonetics []struct {
				Type int    `json:"type"`
				Text string `json:"text"`
			} `json:"phonetics"`
			Explanations []struct {
				Text     string `json:"text"`
				Examples []struct {
					Type      int `json:"type"`
					Sentences []struct {
						Text      string `json:"text"`
						TransText string `json:"trans_text"`
					} `json:"sentences"`
				} `json:"examples"`
				Synonyms []interface{} `json:"synonyms"`
			} `json:"explanations"`
			Relevancys []interface{} `json:"relevancys"`
		} `json:"pos_list"`
	} `json:"words"`
	Phrases  []interface{} `json:"phrases"`
	BaseResp struct {
		StatusCode    int    `json:"status_code"`
		StatusMessage string `json:"status_message"`
	} `json:"base_resp"`
}

func QueryHuoShan(word string, wg *sync.WaitGroup) {
	client := &http.Client{}
	request := HuoShanDictRequest{Text: word, Language: "en"}
	//var data = strings.NewReader(`{"text":"hello","language":"en"}`)
	buf, err := json.Marshal(&request)
	if err != nil {
		log.Fatal(err)
	}
	var data = bytes.NewBuffer(buf)
	req, err := http.NewRequest("POST", "https://translate.volcengine.com/web/dict/match/v1/?msToken=&X-Bogus=DFSzswVYQDcpCDLhSWMzkBt/pL3F&_signature=_02B4Z6wo00001OxXrgwAAIDBZxw17l9jrbDsV6qAAFlovkjJIIMU9EATMS9Tdw-49j7p-K1n4cdGo51cm78Ye2ICZ3L-iQ-U6HgbU1qQezxCn.ikYghdIPtElpIiYafrYfVC-Tuwa6vBPdQ-7a", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("authority", "translate.volcengine.com")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("cookie", "x-jupiter-uuid=16519750659129616; ttcid=783e2ade7358421380c09fb25d87751b33; tt_scid=MScR5VzRv8Ku5jnHj1y2KkEXQnzAZBeBfOXWOGlgu72OlwSpu9cZxsIPmSUU7wUA92ce; s_v_web_id=verify_27f5c06ee18d9523257b3ff1379f47ed; _tea_utm_cache_2018=undefined; i18next=translate")
	req.Header.Set("origin", "https://translate.volcengine.com")
	req.Header.Set("referer", "https://translate.volcengine.com/translate?category=&home_language=zh&source_language=detect&target_language=zh&text=demo")
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="101", "Google Chrome";v="101"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatal("bad StatusCode:", resp.StatusCode, "body", string(bodyText))
	}
	var dictResponse HuoShanDictResponse
	err = json.Unmarshal(bodyText, &dictResponse)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("--------------------火山翻译-----------------------")

	//print phonetics
	Phonetics := dictResponse.Words[0].PosList[0].Phonetics
	fmt.Println(word, "UK:", Phonetics[0].Text, "US:", Phonetics[1].Text)

	//print explanations
	posList := dictResponse.Words[0].PosList
	for i, item := range posList {
		fmt.Printf("解释%d ： %s \n", i+1, item.Explanations[0].Text)
		//for k, example := range item.Examples {
		//	fmt.Println(k, example)
		//}
	}
	fmt.Println("--------------------火山翻译-----------------------")
	wg.Done()
}
