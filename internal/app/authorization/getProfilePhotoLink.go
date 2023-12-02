package authorization

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"main/internal/app/model"
	"net/http"
)

func parseProfilePhotoURL(body io.ReadCloser) (string, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return "", err

	}
	defer body.Close()
	photoURL, _ := doc.Find("#igva_column2_0_avatar").Attr("src")
	return photoURL, nil
}

func (r *Authorization) GetProfilePhotoURL(uid string, client model.ApiRegistration) (string, error) {
	cookies, err := r.GetAuthorizedCookies(uid, client)
	if err != nil {
		return "", err
	}
	req, err := GetAboutInfoRequest(GetCookiesHeader(cookies.Cookie))
	if err != nil {
		return "", err
	}
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return parseProfilePhotoURL(resp.Body)
}
