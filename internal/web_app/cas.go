package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"strings"

	//"gopkg.in/cas.v2"
	"main/internal/database"
	"net/http"
)

type hash map[string]string

func NewCas(service *database.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		CheckAuthorization("dobryninis", "a4c13shj")
	}
}

func CheckAuthorization(login string, password string) (bool, []byte, error) {
	casURL := "https://cas.kai.ru:8443/cas/login"
	postData := url.Values{
		"username":  {login},
		"password":  {password},
		"lt":        {"LT-1677-N5UE5cb3sTos9xdtb4IqH9Ygwx3T1I"},
		"execution": {"e1s1"},
		"_eventId":  {"submit"},
		"Cookie":    {"jsessionid=B744EE55A024CB7FB3813406F5000A78"},
	}
	response, err := http.PostForm(casURL, postData)
	if err != nil {
		panic(err)
	}
	result, _ := io.ReadAll(response.Body)
	fmt.Println(result)
	// Получаем Location из заголовка ответа
	location := response.Header.Get("Location")

	// Проверяем, что браузер перенаправлен на страницу с ticket
	if strings.Contains(location, "ticket=") {
		// Извлекаем ticket из URL-адреса
		ticket := strings.Split(location, "=")[1]

		// Выполняем запрос для проверки билета
		casValidateURL := fmt.Sprintf("https://url.to.cas.server/cas/serviceValidate?ticket=%s&service=%s", ticket, "https://myapp.example.com")
		resp, err := http.Get(casValidateURL)
		if err != nil {
			panic(err)
		}

		// Извлекаем тело ответа
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		// Обрабатываем ответ
		if strings.Contains(string(body), "authenticationSuccess") {
			fmt.Println("Пользователь успешно авторизован.")
		} else {
			fmt.Println("Не удалось авторизовать пользователя.")
		}
	} else {
		fmt.Println("CAS-сервер вернул ошибку: ", location)
	}
	return false, nil, nil
}
