package tools

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const oldKaiURL = "https://old.kai.ru/info/students/brs.php"

type GroupResultAnswer struct {
	Result map[string]string `json:"result"`
}

var Faculties map[int]string

func init() {
	Faculties = make(map[int]string, 6)
	Faculties[1] = "ИАНТЭ"
	Faculties[2] = "ФМФ"
	Faculties[3] = "ИАЭП"
	Faculties[4] = "ИКТЗИ"
	Faculties[5] = "ИРЭТ"
	Faculties[28] = "ВШПИТ и ИИэП"

}

// GetGroupListBRS BRS - Бально рейтинговая система
func GetGroupListBRS(faculty, course string) ([]byte, error) {
	data := fmt.Sprintf("?p_fac=%v&p_kurs=%v", faculty, course)
	resp, err := http.Get(oldKaiURL + data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	resultArray := make(map[string]string)
	doc.Find("select[name=\"p_group\"] option").Each(func(i int, s *goquery.Selection) {
		value, _ := s.Attr("value")
		if value == "" {
			return
		}
		resultArray[s.Text()] = value
	})
	var result GroupResultAnswer
	result.Result = resultArray
	res, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func GetPersonListBRS(faculty, course, group string) ([]byte, error) {
	data := fmt.Sprintf("?p_fac=%v&p_kurs=%v&p_group=%v", faculty, course, group)
	req, err := http.NewRequest("GET", oldKaiURL+data, nil)
	req.Header.Set("Accept-Language", "ru-RU")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	resp, err := http.DefaultClient.Do(req)
	//resp, err := http.Get(oldKaiURL + data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	decoder := charmap.Windows1251.NewDecoder()
	transformReader := transform.NewReader(reader, decoder)

	bufReader := bufio.NewReader(transformReader)

	doc, err := goquery.NewDocumentFromReader(bufReader)
	if err != nil {
		return nil, err
	}
	resultArray := make(map[string]string)
	doc.Find("select[name=\"p_stud\"] option").Each(func(i int, s *goquery.Selection) {
		value, _ := s.Attr("value")
		if value == "" {
			return
		}
		resultArray[value] = s.Text()
	})
	var result GroupResultAnswer
	result.Result = resultArray
	res, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type facResultAnswer struct {
	Result map[int]string `json:"result"`
}

func GetFacultiesListBRS() ([]byte, error) {
	var result facResultAnswer
	result.Result = Faculties
	res, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type ScoreElement struct {
	Index           int    `json:"index"`
	Name            string `json:"name"`
	ScoreCurrent1   int    `json:"scoreCurrent1"`
	ScoreMax1       int    `json:"scoreMax1"`
	ScoreCurrent2   int    `json:"scoreCurrent2"`
	ScoreMax2       int    `json:"scoreMax2"`
	ScoreCurrent3   int    `json:"scoreCurrent3"`
	ScoreMax3       int    `json:"scoreMax3"`
	PreviouslyScore int    `json:"previouslyScore"`
	AdditionalScore int    `json:"additionalScore"`
	Debt            int    `json:"debt"`
	Final           int    `json:"final"`
	Result          string `json:"result"`
}

type ScoreTableAnswer struct {
	Result []ScoreElement `json:"result"`
}

func tryToGetPaidFormStudent(data string) *goquery.Document {
	req, _ := http.NewRequest("POST", oldKaiURL, strings.NewReader(data))
	req.Header.Set("Accept-Language", "ru-RU")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	resp, err := http.DefaultClient.Do(req)
	//resp, err := http.Get(oldKaiURL + data)
	if err != nil {
		log.Printf("Ошибка парсинга: %v", err)
	}
	defer resp.Body.Close()
	reader := bufio.NewReader(resp.Body)
	decoder := charmap.Windows1251.NewDecoder()
	transformReader := transform.NewReader(reader, decoder)
	bufReader := bufio.NewReader(transformReader)

	doc, err := goquery.NewDocumentFromReader(bufReader)
	if err != nil {
		log.Printf("Произошла ошибка парсинга БРС: %v", err)
	}
	return doc
}

func GetScoresStruct(fac, kurs, group, zach, stud int) ([]ScoreElement, error) {
	zachStr := strconv.Itoa(zach)
	if zach < 100000 {
		zachStr = "0" + zachStr
	}
	data := fmt.Sprintf("p_fac=%v&p_kurs=%v&p_group=%v&p_stud=%v&p_zach=%v&p_sub=%v", fac, kurs, group, stud, zachStr, "%CE%F2%EF%F0%E0%E2%E8%F2%FC")
	//decoded, err := url.QueryUnescape("%D0%9F")
	platnikChar := "%CF"
	dataCopy := fmt.Sprintf("p_fac=%v&p_kurs=%v&p_group=%v&p_stud=%v&p_zach=%v&p_sub=%v", fac, kurs, group, stud, platnikChar+zachStr, "%CE%F2%EF%F0%E0%E2%E8%F2%FC")
	req, err := http.NewRequest("POST", oldKaiURL, strings.NewReader(data))
	req.Header.Set("Accept-Language", "ru-RU")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	resp, err := http.DefaultClient.Do(req)
	//resp, err := http.Get(oldKaiURL + data)
	if err != nil {
		log.Printf("Ошибка парсинга: %v", err)
	}
	defer resp.Body.Close()
	reader := bufio.NewReader(resp.Body)
	decoder := charmap.Windows1251.NewDecoder()
	transformReader := transform.NewReader(reader, decoder)
	bufReader := bufio.NewReader(transformReader)

	doc, err := goquery.NewDocumentFromReader(bufReader)
	if err != nil {
		log.Printf("Произошла ошибка парсинга БРС: %v", err)
	}
	table := doc.Find("table[id=\"reyt\"] tr")
	if table.Nodes == nil {
		doc = tryToGetPaidFormStudent(dataCopy) // Ищем данные для платной формы (Прибавляем букву П к номеру)
		table = doc.Find("table[id=\"reyt\"] tr")
		if table.Nodes == nil {
			return []ScoreElement{}, errors.New("Неправильный номер зачетной книжки.")
		}
	}
	scoreElems := make([]ScoreElement, 0)
	table.Find("td").Each(func(i int, s *goquery.Selection) {
		convertToScoreElem(i, s, &scoreElems)
	})
	return scoreElems, nil
}

func GetScores(fac, kurs, group, zach int, stud int) ([]byte, error) {
	scores, err := GetScoresStruct(fac, kurs, group, zach, stud)
	if err != nil {
		return nil, err
	}
	var result ScoreTableAnswer
	result.Result = scores
	res, err := json.Marshal(result)
	if err != nil {
		log.Printf("Ошибка получения баллов БРС: ", err)
		return nil, err
	}
	return res, nil
}

func convertToScoreElem(i int, s *goquery.Selection, scoreElems *[]ScoreElement) {
	if i < 16 { // Пропускаем заголовки таблицы
		return
	}
	index := i - 16
	if index >= 13 {
		index = index % 13
	}
	switch index {
	case 0:
		numberInList, err := strconv.Atoi(s.Text())
		if err != nil {
			log.Printf("Ошибка при генерации списка элементов БРС: \n %v", err)
			return
		}
		newElem := ScoreElement{
			Index: numberInList,
		}
		*scoreElems = append(*scoreElems, newElem)
		return
	case 1: // Имя
		(*scoreElems)[len(*scoreElems)-1].Name = strings.TrimSpace(s.Text())
		return
	case 2: // 1 семестр текущее значение
		(*scoreElems)[len(*scoreElems)-1].ScoreCurrent1, _ = convertToInt(s.Text())
		return
	case 3: // 1 семестр максимальное значение
		(*scoreElems)[len(*scoreElems)-1].ScoreMax1, _ = convertToInt(s.Text())
		return
	case 4:
		(*scoreElems)[len(*scoreElems)-1].ScoreCurrent2, _ = convertToInt(s.Text())
		return
	case 5:
		(*scoreElems)[len(*scoreElems)-1].ScoreMax2, _ = convertToInt(s.Text())
		return
	case 6:
		(*scoreElems)[len(*scoreElems)-1].ScoreCurrent3, _ = convertToInt(s.Text())
		return
	case 7:
		(*scoreElems)[len(*scoreElems)-1].ScoreMax3, _ = convertToInt(s.Text())
		return
	case 8:
		(*scoreElems)[len(*scoreElems)-1].PreviouslyScore, _ = convertToInt(s.Text())
		return
	case 9:
		(*scoreElems)[len(*scoreElems)-1].AdditionalScore, _ = convertToInt(s.Text())
		return
	case 10:
		(*scoreElems)[len(*scoreElems)-1].Debt, _ = convertToInt(s.Text())
		return
	case 11:
		(*scoreElems)[len(*scoreElems)-1].Final, _ = convertToInt(s.Text())
		return
	case 12:
		(*scoreElems)[len(*scoreElems)-1].Result = strings.TrimSpace(s.Text())
		return
	case 13:
		(*scoreElems)[len(*scoreElems)-1].Name = strings.TrimSpace(s.Text())
		return

	}
}

func convertToInt(str string) (int, error) {
	if str == "" {
		return 0, nil
	}
	return strconv.Atoi(str)
}
