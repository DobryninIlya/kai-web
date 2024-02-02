package get_request

import (
	"encoding/json"
	"io"
	"log"
	"main/internal/app/model"
	"net/http"
	"strconv"
)

const link = "https://kai.ru/"
const parametersGroup = "raspisanie?p_p_id=pubStudentSchedule_WAR_publicStudentSchedule10&p_p_lifecycle=2&p_p_resource_id=getGroupsURL"
const parametersSchedule = "raspisanie?p_p_id=pubStudentSchedule_WAR_publicStudentSchedule10&p_p_lifecycle=2&p_p_state=normal&p_p_mode=view&p_p_resource_id=schedule&groupId="

func getGroupsStruct(body []byte) []GroupInfo {
	var list []GroupInfo
	err := json.Unmarshal(body, &list)
	if err != nil {
		log.Println(err)
	}
	return list
}

func GetGroupsList() []GroupInfo {
	resp, err := http.Get(link + parametersGroup)
	if err != nil {
		log.Println("Ошибка получения списка групп \n")
		log.Println(link + parametersGroup)
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return getGroupsStruct(body)

}

func getScheduleStruct(body []byte) model.Schedule {
	var shed model.Schedule
	err := json.Unmarshal(body, &shed)
	if err != nil {
		log.Println(err)
	}
	return shed
}

func GetScheduleByGroup(group int) model.Schedule {
	resp, err := http.Get(link + parametersSchedule + strconv.Itoa(group))
	if err != nil {
		log.Println("Ошибка получения расписания группы \n")
		log.Println(link + parametersGroup + strconv.Itoa(group))
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return getScheduleStruct(body)

}

func GetUnmarshaledSchedule(data []byte) model.Schedule {
	var result model.Schedule
	err := json.Unmarshal(data, &result)
	if data == nil {
		return model.Schedule{}
	}
	if err != nil {
		log.Printf("Ошибка анмаршалинга #{err}")
		return model.Schedule{}
	}
	return result
}

func GetMarshaledSchedule(data model.Schedule) []byte {
	result, err := json.Marshal(data)
	if err != nil {
		log.Printf("Ошибка маршалинга #{err}")
	}
	return result
}
