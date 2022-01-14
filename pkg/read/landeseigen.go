package read

import (
	"encoding/json"
	"immo/pkg/telegram"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var knownURLs []string

func Landeseigen() {
	endpoint := "https://inberlinwohnen.de/wp-content/themes/ibw/skript/wohnungsfinder.php"
	data := url.Values{}
	data.Set("q", "wf-save-srch")
	data.Set("miete_max", "550")
	data.Set("rooms_min", "2")
	data.Set("seniorenwohnung", "false")
	data.Set("bez[]", "02_00")
	data.Set("bez[]", "11_00")
	data.Set("wbs", "0")

	client := &http.Client{}
	r, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode())) // URL-encoded payload
	if err != nil {
		log.Fatal(err)
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := client.Do(r)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var dataResp map[string]string
	err = json.Unmarshal(body, &dataResp)

	split := strings.Split(dataResp["searchresults"], "href=")

	var urls []string
	for i := range split {
		urls = append(urls, GetStringInBetweenTwoString(split[i], "\"", ".html"))
	}

	var found bool

	if dataResp["headline"] != "Wir haben 0 Wohnungen für Sie gefunden" {
		for i := range urls {
			found = false
			if urls[i] != "" {
				parse, _ := url.Parse("https://inberlinwohnen.de/wohnungsfinder" + urls[i] + ".html")
				for x := range knownURLs {
					if knownURLs[x] == parse.String() {
						found = true
					}
				}
				if found == false {
					knownURLs = append(knownURLs, parse.String())
					telegram.SendTextToTelegramChat(1167392515, parse.String())
				}
			}
		}
	}
}

func GetStringInBetweenTwoString(str string, startS string, endS string) (result string) {
	s := strings.Index(str, startS)
	if s == -1 {
		return result
	}
	newS := str[s+len(startS):]
	e := strings.Index(newS, endS)
	if e == -1 {
		return result
	}
	result = newS[:e]
	return result
}
