package authorization

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"main/internal/app/model"
	"net/http"
	"strconv"
)

func GetAttestationRequest(cookies string) (*http.Request, error) {
	req, err := http.NewRequest("GET", AttestationPage, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Cookie", cookies)
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	return req, nil
}

func parseAttestationTable(body io.ReadCloser) ([]model.Discipline, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	tables := doc.Find("table")
	tableCount := tables.Length()

	if tableCount < 3 {
		log.Println("table attestation not found")
		return nil, ErrAttestationNotFound
	}

	lastTable := tables.Eq(tableCount - 3)

	var disciplines []model.Discipline

	lastTable.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
		discipline := model.Discipline{}

		row.Find("td").Each(func(j int, cell *goquery.Selection) {
			switch j {
			case 0:
				discipline.Number, _ = strconv.Atoi(cell.Text())
			case 1:
				discipline.Name = cell.Find("a").Text()
			case 2, 3, 4, 5, 6, 7, 8, 9, 10, 11:
				assessment := model.Assessment{}
				assessment.YourScore, _ = strconv.Atoi(cell.Text())
				assessment.MaxScore, _ = strconv.Atoi(cell.Next().Text())
				discipline.Assessments = append(discipline.Assessments, assessment)
			case 12:
				discipline.PreliminaryGrade = cell.Text()
			case 13:
				discipline.AdditionalPoints, _ = strconv.Atoi(cell.Text())
			case 14:
				discipline.Debts, _ = strconv.Atoi(cell.Text())
			case 15:
				discipline.FinalGrade, _ = strconv.Atoi(cell.Text())
			case 16:
				discipline.TraditionalGrade = cell.Text()
			}
		})

		disciplines = append(disciplines, discipline)
	})

	return disciplines, nil
}

// GetAttestationList возвращает список с баллами БРС (либо ошибку)
// Важно отметить, что возвращает только последний семестр последнего (текущего) курса (группы)
func (r *Authorization) GetAttestationList(uid string, client model.ApiRegistration) ([]model.Discipline, error) {
	value, ok := r.GetAttestations(uid)
	if ok && value != nil {
		return value, nil
	}
	cookies, err := r.GetAuthorizedCookies(uid, client)
	if err != nil {
		return nil, err
	}
	req, err := GetAttestationRequest(GetCookiesHeader(cookies.Cookie))
	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result, err := parseAttestationTable(resp.Body)
	if err != nil {
		return nil, err
	}
	r.SetAttestations(uid, result)
	return result, nil
}
