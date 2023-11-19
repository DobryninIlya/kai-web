package authorization

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"main/internal/app/model"
	"net/http"
	"strconv"
	"strings"
)

func GetGroupRequest(cookies string) (*http.Request, error) {
	req, err := http.NewRequest("GET", MyGroupPage, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Cookie", cookies)
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	return req, nil
}

func (r *Authorization) GetGroupNum(uid string, apiClient model.ApiRegistration) (int, error) {
	cookies, err := r.GetAuthorizedCookies(uid, apiClient)
	if err != nil {
		return 0, err
	}
	req, err := GetGroupRequest(GetCookiesHeader(cookies.Cookie))
	if err != nil {
		return 0, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	return ParseGroupNum(resp.Body)
}

func ParseGroupNum(body io.ReadCloser) (int, error) {
	defer body.Close()
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return 0, err
	}
	defer body.Close()
	input := doc.Find("input.field[name='_myGroup_WAR_myGroup10_orgUnit']")
	groupElement := input.Last()
	group := groupElement.Parent().Text()
	group = strings.ReplaceAll(group, "\n", "")
	group = strings.ReplaceAll(group, "\t", "")
	groupNum, err := strconv.Atoi(group)
	if err != nil {
		return 0, err
	}
	return groupNum, nil
}
