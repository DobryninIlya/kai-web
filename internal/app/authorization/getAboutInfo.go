package authorization

import (
	"github.com/PuerkitoBio/goquery"
	"io"

	"main/internal/app/model"
	"net/http"
)

func GetAboutInfoRequest(cookies string) (*http.Request, error) {
	req, err := http.NewRequest("GET", AboutPage, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Cookie", cookies)
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	return req, nil
}

func ParseAboutInfo(body io.ReadCloser) (model.SiteUserInfo, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return model.SiteUserInfo{}, err
	}
	defer body.Close()
	user := model.SiteUserInfo{
		FirstName:  doc.Find("#_aboutMe_WAR_aboutMe10_firstName").AttrOr("value", ""),
		LastName:   doc.Find("#_aboutMe_WAR_aboutMe10_lastName").AttrOr("value", ""),
		MiddleName: doc.Find("#_aboutMe_WAR_aboutMe10_middleName").AttrOr("value", ""),
	}
	return user, nil
}

func (r *Authorization) GetAboutInfo(uid string, client model.ApiRegistration) (model.SiteUserInfo, error) {
	cookies, err := r.GetAuthorizedCookies(uid, client)
	if err != nil {
		return model.SiteUserInfo{}, err
	}
	req, err := GetAboutInfoRequest(GetCookiesHeader(cookies.Cookie))
	if err != nil {
		return model.SiteUserInfo{}, err
	}
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return model.SiteUserInfo{}, err
	}
	defer resp.Body.Close()
	return ParseAboutInfo(resp.Body)
}
