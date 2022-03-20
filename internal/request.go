package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func fillParam() []byte {
	type ParamsGetDevices struct {
		Action     string `json:"action"`
		Page       string `json:"page"`
		From       int    `json:"from"`
		To         int    `json:"to"`
		City       int    `json:"city"`
		Type       string `json:"type"`
		Give       string `json:"give"`
		Get        string `json:"get"`
		Commission int    `json:"commission"`
		Sort       string `json:"sort"`
		Range      string `json:"range"`
		Sortm      int    `json:"sortm"`
		Tsid       int    `json:"tsid"`
	}

	params := ParamsGetDevices{
		"getrates",
		"rates",
		59,
		10,
		0,
		"",
		"",
		"",
		0,
		"from",
		"asc",
		0,
		0,
	}

	strParams, _ := json.Marshal(params)
	return strParams
}

func toUTF(win string) []byte {
	sr := strings.NewReader(win)
	tr := transform.NewReader(sr, charmap.Windows1251.NewDecoder())
	buf, err := ioutil.ReadAll(tr)
	if err != err {
		log.Println(err)
		return nil
	}
	return buf
}

type Exchanges struct {
	Name   string
	Url    string
	Value  float64
	Volume uint64
}

func (m *MetricsStructure) GetRequest(mName, url string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("User-Agent", " Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36")
	res, err := client.Do(req)
	if err != nil {
		return
	}

	str, _ := ioutil.ReadAll(res.Body)
	b := toUTF(string(str))

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}

	res.Body.Close()

	// Find the review items
	doc.Find("table#content_table tr").Each(func(i int, s *goquery.Selection) {
		//if i > 0 {
		item := Exchanges{}
		item.Url, _ = s.Find("td").Eq(1).Find("div").Eq(0).Find("a").Attr("href")
		item.Name = s.Find("td").Eq(1).Find("div").Eq(0).Find("div").Eq(0).Text()

		var value string
		if mName == "exchangesBuy" {
			value := s.Find("td").Eq(2).Find("div").Eq(0).Text()
			small := s.Find("td").Eq(2).Find("div").Eq(0).Find("small").Eq(0).Text()
			value = strings.TrimSpace(strings.ReplaceAll(value, small, ""))
			item.Value, _ = strconv.ParseFloat(value, 10)
		} else {
			value := s.Find("td").Eq(3).Text()
			small := s.Find("td").Eq(3).Find("small").Eq(0).Text()
			value = strings.TrimSpace(strings.ReplaceAll(value, small, ""))
			item.Value, _ = strconv.ParseFloat(value, 10)
		}

		value = s.Find("td").Eq(4).Text()
		value = strings.TrimSpace(strings.ReplaceAll(value, " ", ""))
		item.Volume, _ = strconv.ParseUint(value, 10, 64)

		if item.Name != "" {

			m.MChannel <- Metric{fmt.Sprintf("%s{instance=\"%s\"}", mName, item.Name), item.Value}

		}
		//	fmt.Printf("Review %d: %s\n", i, title)
		//}
	})

	//doc.Find("table#content_table").Each(func(i int, s *goquery.Selection) {
	//	// For each item found, get the title
	//	title := s.Find("a").Text()
	//	fmt.Printf("Review %d: %s\n", i, title)
	//})

}

func PostRequest() {
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://www.bestchange.ru/action.php", bytes.NewReader(fillParam()))
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("Content-Length", strconv.Itoa(len(fillParam())))
	req.Header.Add("Content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Cookie", "PHPSESSID=prnpq4rah41elboad0f8rhuv41; userid=222f8224b2faf0655aeb2211569a9897; source=o%B53%28%B8%CE%F9%8E%87%E5%CF%82%DA%E6%F4%9F%AE%0E%81%C5%40rY%7C%80%CA%C5%3F%C3%E1%FF%DB%0FCU%5D9%0B%A2z%21%AC%2F%D3%B6%E1%1D%CEp%B2%1F%EB%83; pixel=1; _ga=GA1.2.236511636.1647775515; _gid=GA1.2.625572206.1647775515; _ym_uid=1647775515104562544; _ym_d=1647775515; _ym_isad=1; history=1647775393-749-1-59-10-112.61218000-24400140.00a1647777090-638-1-59-10-114.00000000-14000.00")
	req.Header.Add("Host", "www.bestchange.ru")
	req.Header.Add("Origin", "https://www.bestchange.ru")
	req.Header.Add("Referer", "https://www.bestchange.ru/visa-mastercard-rub-to-tether-trc20.html")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("User-Agent", " Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36")

	req.Header.Add("sec-ch-ua", `"Not A;Brand";v="99", "Chromium";v="98", "Google Chrome";v="98"`)
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "Linux")

	res, err := client.Do(req)

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(b))
}
