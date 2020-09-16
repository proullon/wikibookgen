package editor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

type LocateWikiFileResponsePage struct {
	Pageid          int    `json:"pageid"`
	Ns              int    `json:"ns"`
	Title           string `json:"title"`
	Imagerepository string `json:"imagerepository"`
	Imageinfo       []struct {
		URL                 string `json:"url"`
		Descriptionurl      string `json:"descriptionurl"`
		Descriptionshorturl string `json:"descriptionshorturl"`
	} `json:"imageinfo"`
}

type LocateWikiFileResponse struct {
	Batchcomplete string `json:"batchcomplete"`
	Query         struct {
		Normalized []struct {
			From string `json:"from"`
			To   string `json:"to"`
		} `json:"normalized"`
		Pages map[string]LocateWikiFileResponsePage `json:"pages"`
	} `json:"query"`
}

func LocateWikiFile(s string, lang string) (string, error) {

	url, err := locatewikifile(s, "commons.wikimedia.org")
	if err == nil {
		return url, nil
	}

	return locatewikifile(s, fmt.Sprintf("%s.wikipedia.org", lang))
}

func locatewikifile(s string, domain string) (string, error) {
	s = strings.TrimSpace(s)
	s = strings.Replace(s, " ", "_", -1)

	urlfmt := `https://%s/w/api.php?action=query&titles=File%%3A%s&prop=imageinfo&iiprop=url&format=json`
	url := fmt.Sprintf(urlfmt, domain, s)

	log.Debugf("Querying %s", url)

	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return "", err
	}

	r := LocateWikiFileResponse{}
	err = json.Unmarshal(data, &r)
	if err != nil {
		return "", fmt.Errorf("unmarshal: %s", err)
	}

	for _, v := range r.Query.Pages {
		if len(v.Imageinfo) == 0 {
			return "", fmt.Errorf("not found")
		}

		return v.Imageinfo[0].URL, nil
	}
	return "", fmt.Errorf("not found")
}
