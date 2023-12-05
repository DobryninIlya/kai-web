package authorization

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"io"
	"main/internal/app/model"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
)

const (
	KaiURL = "https://kai.ru"
)

func parseFormActionURL(body io.ReadCloser) (string, map[string]string, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return "", nil, err

	}
	defer body.Close()
	photoURL, _ := doc.Find("#_2_fm").Attr("action")
	inputs := doc.Find("form input")
	inputMap := make(map[string]string, 3)
	// Проходимся по каждому найденному элементу и выводим его атрибуты
	inputs.Each(func(index int, input *goquery.Selection) {
		name, exists := input.Attr("name")
		if exists {
			value, _ := input.Attr("value")
			inputMap[name] = value
		}
	})
	return photoURL, inputMap, nil
}

//	func parseDownloadFormURL(body io.ReadCloser) (string, error) {
//		doc, err := goquery.NewDocumentFromReader(body)
//		if err != nil {
//			return "", err
//
//		}
//		defer body.Close()
//		photoURL, _ := doc.Find("#_aboutMe_WAR_aboutMe10_changeLogo_iframe_").Attr("src")
//		return photoURL, nil
//	}
func parseDownloadFormURL(body io.ReadCloser) (string, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return "", err

	}
	defer body.Close()
	var scriptText string
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "Liferay.LogoSelector") {
			scriptText = s.Text()
			return
		}
	})

	// найти значение editLogoURL
	editLogoURL := ""
	if strings.Contains(scriptText, "editLogoURL") {
		startIndex := strings.Index(scriptText, "editLogoURL: '") + len("editLogoURL: '")
		endIndex := strings.Index(scriptText[startIndex:], "'")
		editLogoURL = scriptText[startIndex : startIndex+endIndex]
	}

	fmt.Println(editLogoURL)
	return editLogoURL, nil
}

func (r *Authorization) GetDownloadFormURL(uid string, client model.ApiRegistration) (string, error) {
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
	return parseDownloadFormURL(resp.Body)
}

func GetFormRequest(url string, cookies string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("cookie", cookies)
	return req, nil
}

func (r *Authorization) GetFormActionURL(uid string, client model.ApiRegistration, url string) (string, map[string]string, error) {
	cookies, err := r.GetAuthorizedCookies(uid, client)
	if err != nil {
		return "", nil, err
	}
	req, err := GetFormRequest(url, GetCookiesHeader(cookies.Cookie))
	if err != nil {
		return "", nil, err
	}
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()
	return parseFormActionURL(resp.Body)
}

func GetUploadPhotoRequest(url string, cookies string, inputs map[string]string, file io.Reader, cmd string) (*http.Request, error) {
	// Создаем буфер для записи multipart/form-data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	// Добавляем файл к форме
	part, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Disposition": []string{`form-data; name="_2_fileName"; filename="success.png"`},
		"Content-Type":        []string{"image/png"},
	})
	if err != nil {
		return nil, err
	}
	io.Copy(part, file)
	// Добавляем остальные поля
	for key, value := range inputs {
		writer.WriteField(key, value)
	}
	if cmd != "" {
		writer.WriteField("_2_cmd", cmd)
	}
	//writer.WriteField("_2_p_u_i_d_", "add_temp")
	writer.Close()
	// Создаем запрос с телом формы
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return nil, err
	}
	// Устанавливаем заголовки
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Cookie", cookies)
	return req, nil
}

func (r *Authorization) UploadProfilePhoto(uid string, client model.ApiRegistration, file io.Reader, cmd string) error {
	cookies, err := r.GetAuthorizedCookies(uid, client)
	if err != nil {
		return err
	}
	url, err := r.GetDownloadFormURL(uid, client)
	if err != nil {
		r.log.Logf(logrus.ErrorLevel, "Error while getting download form url: %s", err.Error())
		return err
	}
	formActionURL, inputs, err := r.GetFormActionURL(uid, client, url)
	if err != nil {
		r.log.Logf(logrus.ErrorLevel, "Error while getting form action url: %s", err.Error())
		return err
	}
	req, err := GetUploadPhotoRequest(formActionURL, GetCookiesHeader(cookies.Cookie), inputs, file, cmd)
	if err != nil {
		r.log.Logf(logrus.ErrorLevel, "Error while getting upload photo request: %s", err.Error())
		return err
	}
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		r.log.Logf(logrus.ErrorLevel, "Error while uploading photo: %s", err.Error())
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		r.log.Logf(logrus.ErrorLevel, "Error while uploading photo: %s", err.Error())
		return err
	}
	return nil
}

func (r *Authorization) ChangeProfilePhoto(uid string, client model.ApiRegistration, file io.Reader) error {
	r.UploadProfilePhoto(uid, client, file, "add_temp") // Добавляем превью фотографии
	r.UploadProfilePhoto(uid, client, file, "")         // Загружаем ее же повторно (таковы требования сайта)
	return nil
}
