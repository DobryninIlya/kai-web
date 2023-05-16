package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Service struct {
	//conn *pgxpool.Conn
	conn *pgxpool.Pool
}

type RegistrationData struct {
	Role          int    `json:"id"`
	Identificator string `json:"Identificator"`
	VkId          string `json:"vk_id"`
}

type GroupInfo struct {
	Id    int    `json:"id"`
	Group string `json:"group"`
	Form  string `json:"forma,omitempty"`
}

const examUrl = "https://kai.ru/raspisanie?p_p_id=pubStudentSchedule_WAR_publicStudentSchedule10&p_p_lifecycle=2&p_p_resource_id=examSchedule&groupId="

var chetn int
var reservedDict = []string{"чет", "неч", "чет/неч", "неч/чет"}

// var groupList = make([]GroupInfo, 0)
var groupList []GroupInfo

func init() {
	chetnStr, _ := strconv.Atoi(os.Getenv("CHETN"))
	chetn = chetnStr

}

func (s *Service) IsRegistredUser(id int) bool {

	rows, err := s.conn.Query(context.Background(), "SELECT * FROM users WHERE id_vk = $1", id)
	if err != nil {
		log.Println("Ошибка запроса группы")
		log.Println(err)
	}
	defer rows.Close()
	rows.Next()
	val, err := rows.Values()
	if val == nil {
		return false
	}
	return true
}

type Prepod struct {
	lessonType []string
	name       string
	lesson     string
}

func (s Service) GetTeacherListStruct(id int) []Prepod {
	groupId := s.GetGroupByUserId(id)
	sched := s.GetScheduleByGroup(groupId)
	prepodList := make([]Prepod, 0)
	v := reflect.ValueOf(sched)
	for i := 0; i < v.NumField(); i++ { // перебираем все поля структуры
		field := v.Field(i)
		if field.Kind() == reflect.Slice { // проверяем, что поле является срезом
			sliceValue := field.Interface()             // получаем значение среза
			if slice, ok := sliceValue.([]Lesson); ok { // проверяем, что значение среза имеет тип []Lesson
				for _, lesson := range slice {
					lesson.PrepodName = strings.TrimSpace(lesson.PrepodName)
					lesson.DisciplType = strings.TrimSpace(lesson.DisciplType)
					lesson.DisciplName = strings.TrimSpace(lesson.DisciplName)
					fmt.Printf("%s | %s \n", lesson.DisciplType, lesson.PrepodName) // выводим строку "ID === Name" для каждого элемента массива
					added := false
					for k, prepod := range prepodList {
						if prepod.name == lesson.PrepodName && prepod.lesson == lesson.DisciplName {
							if !CheckInSlice(prepod.lessonType, lesson.DisciplType) {
								prepod.lessonType = append(prepod.lessonType, lesson.DisciplType)
								prepodList[k] = prepod
							}
							added = true
							break

						}
					}
					if added {
						continue
					}
					newPrepod := Prepod{
						lessonType: make([]string, 1),
						name:       lesson.PrepodName,
						lesson:     lesson.DisciplName,
					}
					newPrepod.lessonType[0] = lesson.DisciplType
					prepodList = append(prepodList, newPrepod)
				}
			}
		}
	}

	return prepodList

}

func (s Service) GetExamListStruct(id int) []ExamStruct {
	groupId := s.GetGroupByUserId(id)
	resp, err := http.Get(examUrl + strconv.Itoa(groupId))
	if err != nil {
		log.Printf("Ошибка зпроса расписания экзаменов в database: %v", err.Error())
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	examStruct, err := GetExamStruct([]byte(body))
	if err != nil {
		return nil
	}
	return examStruct

}

func (s *Service) GetGroupByUserId(id int) int {
	rows, err := s.conn.Query(context.Background(), "SELECT groupp FROM users WHERE id_vk = $1", id)
	if err != nil {
		log.Println("Ошибка запроса группы")
		log.Println(err)
	}
	var result int
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&result)
		if err != nil {
			return 0
		}
		return result
	}
	if err := rows.Err(); err != nil {
		// обработка ошибки
		fmt.Println("Ошибка:", err)
		return 0
	}

	if !rows.Next() {
		// запрос вернул пустой результат
		fmt.Println("Результат запроса пустой")
		return 0
	}
	return result
}

func (s *Service) GetScheduleByGroup(id int) Schedule {
	rows, err := s.conn.Query(context.Background(), "SELECT shedule FROM saved_timetable WHERE groupp = $1", id)
	if err != nil {
		log.Println("Ошибка запроса расписания GetScheduleByGroup")
		log.Println(err)
	}
	var schedule string
	defer rows.Close()
	rows.Next()
	rows.Scan(&schedule)
	scheduleStruct := GetScheduleStruct([]byte(schedule))
	return scheduleStruct
}

func (s *Service) GetIdByGroup(id int) int {
	idStr := strconv.Itoa(id)
	if len(groupList) == 0 {
		rows, err := s.conn.Query(context.Background(), "SELECT shedule FROM saved_timetable WHERE groupp=$1", 1)

		if err != nil {
			log.Println("Ошибка получения списка групп")
			log.Println(err)
			return 0
		}
		defer rows.Close()
		var result string
		for rows.Next() {
			rows.Scan(&result)
		}
		json.Unmarshal([]byte(result), &groupList)
		if err != nil {
			log.Println(err)
		}
	}
	for _, group := range groupList {
		if group.Group == idStr {
			return group.Id
		}
	}
	return 0
}

func (s *Service) MakeRegistration(data RegistrationData) (bool, error) {
	var groupId int
	var login string
	groupReal, err := strconv.Atoi(data.Identificator)
	if err == nil {
		login = ""
		groupId = s.GetIdByGroup(groupReal)
		if groupId == 0 {
			return false, errors.New(fmt.Sprintf("Группы %v не существует.", groupReal))
		}
	} else {
		login = data.Identificator
	}

	row, err := s.conn.Query(context.Background(), "INSERT INTO public.users(id_vk, name, groupp, distribution, admlevel, groupreal, \"dateChange\", balance, distr, warn, expiration, banhistory, ischeked, role, login, potok_lecture, has_own_shed, affiliate)"+
		"VALUES ($1, '', $2, 1, 1, $3, Now(), 0, 0, 0, '2020-01-01', 0, 0, $4, $5, true, false, false)",
		data.VkId, groupId, groupReal, data.Role, login)
	row.Next()
	row.Close()
	if err != nil {
		log.Println("Ошибка регистрации. INSERT")
		log.Println(err)
		return false, errors.New("Ошибка регистрации. Попробуйте позже")
	}
	return true, nil
}

func (s *Service) GetCurrentDaySchedule(uId int, margin int) ([]Lesson, time.Time) {
	day := time.Now().AddDate(0, 0, margin)
	dayNum := day.Weekday()

	groupSchedule := s.GetScheduleByGroup(s.GetGroupByUserId(uId))
	lessons := make([]Lesson, 0, 4)
	switch {
	case dayNum == 1:
		lessons = groupSchedule.Day1
	case dayNum == 2:
		lessons = groupSchedule.Day2
	case dayNum == 3:
		lessons = groupSchedule.Day3
	case dayNum == 4:
		lessons = groupSchedule.Day4
	case dayNum == 5:
		lessons = groupSchedule.Day5
	case dayNum == 6:
		lessons = groupSchedule.Day6
	case dayNum == 7:
		lessons = []Lesson{}
	}

	return formScheduleList(lessons, margin), day
}

func formScheduleList(lessons []Lesson, margin int) []Lesson {
	_, week := time.Now().ISOWeek()
	result := make([]Lesson, 0)
	isEven := (week%2 + chetn) == 0
	for _, lesson := range lessons {
		date := strings.TrimSpace(lesson.DayDate)
		re := regexp.MustCompile(`^[.-]+[.\s]+$`) // регулярное выражение для точек и тире
		if re.MatchString(date) {
			date = ""
		}
		if re.MatchString(lesson.AudNum) {
			lesson.AudNum = ""
		}
		if date == "чет" && isEven {
			result = append(result, lesson)
		} else if date == "неч" && !isEven {
			result = append(result, lesson)
		} else if (date == "чет/неч" && isEven) || (date == "неч/чет" && !isEven) {
			lesson.DayDate = "[1 гр.]"
			result = append(result, lesson)
		} else if (date == "неч/чет" && isEven) || (date == "чет/неч" && !isEven) {
			lesson.DayDate = "[2 гр.]"
			result = append(result, lesson)
		} else if ex1, ex2 := isContainDate(date, margin); ex1+ex2 != "" {
			lesson.DayDate = getSubgroupForDate(date, ex1, ex2, isEven)
			result = append(result, lesson)
		} else if !isContainsInDict(date, margin) {
			result = append(result, lesson)
		}
	}
	return result
}

func isContainsInDict(date string, margin int) bool {
	for _, v := range reservedDict {
		if v == date {
			return true
		}
	}

	return false
}

func isContainDate(data string, margin int) (string, string) {
	date := time.Now().AddDate(0, 0, margin)
	day := date.Day()
	dayString := strconv.Itoa(day)

	month := int(date.Month())
	monthString := strconv.Itoa(month)

	ex1 := dayString + "." + monthString
	if day < 0 {
		dayString = "0" + dayString
	}
	if month < 0 {
		monthString = "0" + monthString
	}
	ex2 := dayString + "." + monthString
	if strings.Contains(data, ex1) || strings.Contains(data, ex2) {
		return ex1, ex2
	}
	return "", ""
}

func getSubgroupForDate(data, ex1, ex2 string, isEven bool) string {
	if strings.Contains(data, "/") {
		parts := strings.Split(data, "/")
		if len(parts) != 2 {
			return data
		}
		if strings.Contains(parts[0], ex1) || strings.Contains(parts[0], ex2) {
			if isEven {
				return "[1 гр.]"
			} else {
				return "[2 гр.]"
			}
		} else if strings.Contains(parts[1], ex1) || strings.Contains(parts[0], ex2) {
			if isEven {
				return "[2 гр.]"
			} else {
				return "[1 гр.]"
			}
		}
	} else {
		if strings.Contains(data, ex1) || strings.Contains(data, ex2) {
			return ex1
		}

	}
	return ""
}
func NewService(ctx context.Context, pgConfig pgxpool.Config) *Service {

	config, _ := pgxpool.ParseConfig("")
	config.ConnConfig.Host = pgConfig.ConnConfig.Host
	config.ConnConfig.Port = pgConfig.ConnConfig.Port
	config.ConnConfig.User = pgConfig.ConnConfig.User
	config.ConnConfig.Password = pgConfig.ConnConfig.Password
	config.ConnConfig.Database = pgConfig.ConnConfig.Database
	conn, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	return &Service{
		conn: conn,
	}
}
