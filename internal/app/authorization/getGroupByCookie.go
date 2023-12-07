package authorization

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"main/internal/app/model"
	"net/http"
	"strconv"
	"strings"
)

func GetGroupRequest(cookies string) (*http.Request, error) {
	req, err := http.NewRequest("GET", MyJobForm, nil)
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

func ParseGroupNum(body io.Reader) (int, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return 0, err
	}

	// Find the row with the group information
	row := doc.Find("tr:contains('Группа №')")
	if row.Length() == 0 {
		return 0, fmt.Errorf("Group information not found")
	}

	// Extract the group number
	groupNum := row.Find("td.new_white_td1").Text()
	groupNum = strings.TrimSpace(groupNum)
	groupNumInt, err := strconv.Atoi(groupNum)
	if err != nil {
		return 0, err
	}

	return groupNumInt, nil
}
