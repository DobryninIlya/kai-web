package parser

import (
	"encoding/json"
	"io"
	"log"
	"main/internal/app/model"
	"net/http"
	"strconv"
)

const examUrl = "https://kai.ru/raspisanie?p_p_id=pubStudentSchedule_WAR_publicStudentSchedule10&p_p_lifecycle=2&p_p_resource_id=examSchedule&groupId="

func GetExamListStruct(groupId int) ([]model.Exam, error) {
	resp, err := http.Get(examUrl + strconv.Itoa(groupId))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	examStruct, err := GetExamStruct(body)
	if err != nil {
		return nil, err
	}
	return examStruct, nil

}
func GetExamStruct(body []byte) ([]model.Exam, error) {
	var shed []model.Exam
	err := json.Unmarshal(body, &shed)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return shed, nil
}
