package sqlstore

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"main/internal/app/database"
	"main/internal/app/model"
	"main/internal/app/store"
	_ "main/internal/app/store"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var chetn int
var reservedDict = []string{"чет", "неч", "чет/неч", "неч/чет"}

// UserRepository ...
type ScheduleRepository struct {
	store *Store
}

type GroupInfo struct {
	Id    int    `json:"id"`
	Group string `json:"group"`
	Form  string `json:"forma,omitempty"`
}

var groupList []GroupInfo

func (r ScheduleRepository) GetIdByGroup(id int) (int, error) {
	var result string
	if err := r.store.db.QueryRow(
		"SELECT shedule FROM saved_timetable WHERE groupp=$1",
		1,
	).Scan(&result); err != nil {
		if err == sql.ErrNoRows {
			return 0, store.ErrRecordNotFound
		}

		return 0, err
	}
	json.Unmarshal([]byte(result), &groupList)
	idStr := strconv.Itoa(id)
	for _, group := range groupList {
		if group.Group == idStr {
			return group.Id, nil
		}
	}

	return 0, nil
}

func (r ScheduleRepository) MarkDeletedLesson(user model.User, lessonId int, uniqString string) (int, error) {
	var result int
	if err := r.store.db.QueryRow(
		"INSERT INTO public.deleted_lessons(groupid, creator, creator_platform, lesson_id, date, uniqstring) VALUES ($1, $2, $3, $4, NOW(), $5) RETURNING id;",
		user.Group,
		user.ID,
		"vk",
		lessonId,
		uniqString,
	).Scan(&result); err != nil {
		return 0, err
	}
	return result, nil
}

func (r ScheduleRepository) ReturnDeletedLesson(user model.User, lessonId int, uniqString string) (int, error) {
	var result int
	if err := r.store.db.QueryRow(
		"DELETE FROM public.deleted_lessons WHERE lesson_id=$1 AND uniqstring=$2 RETURNING id;",

		lessonId,
		uniqString,
	).Scan(&result); err != nil {
		return 0, err
	}
	return result, nil
}

func (r ScheduleRepository) GetDeletedLessonsByGroup(groupId int) ([]model.DeletedLessonsMin, error) {
	var resultStructList []model.DeletedLessonsMin
	rows, err := r.store.db.Query(
		"SELECT lesson_id, uniqstring FROM public.deleted_lessons WHERE groupid=$1;",
		groupId,
	)
	if err != nil {
		return nil, err
	}
	resultStructList = make([]model.DeletedLessonsMin, 0)
	for rows.Next() {
		var lessonId model.DeletedLessonsMin
		err := rows.Scan(&lessonId.LessonId, &lessonId.Uniqstring)
		if err != nil {
			log.Printf("Ошибка сканирования в GetDeletedLessonsByGroup: %v", err)
			return nil, err
		}
		resultStructList = append(resultStructList, lessonId)
	}
	return resultStructList, nil
}

func (r ScheduleRepository) NewLesson(user model.User, lesson model.LessonNew) (int, error) {
	var result int
	jsonData, err := json.Marshal(lesson)
	if err != nil {
		return 0, errors.New("ошибка создания json структуры в NewLesson")
	}
	if err := r.store.db.QueryRow(
		"INSERT INTO public.lessons(groupid, daynum, lesson_data) VALUES ($1, $2, $3) RETURNING id;",
		user.Group,
		lesson.DayNum,
		jsonData,
	).Scan(&result); err != nil {
		return 0, err
	}
	return result, nil
}

func formScheduleList(lessons []model.Lesson, margin int) []model.Lesson {
	_, week := time.Now().AddDate(0, 0, margin).ISOWeek()
	result := make([]model.Lesson, 0)
	isEven := (week%2 + chetn) == 0
	for _, lesson := range lessons {
		date := strings.TrimSpace(lesson.DayDate)
		date = strings.ToLower(date)
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
			lesson.DayDate = getSubgroupForDate(date, ex1, ex2)
			result = append(result, lesson)
		} else if !isContainsInDict(date) {
			result = append(result, lesson)
		}
	}
	return result
}

// Содержит ли тип пары дату
func isContainDate(data string, margin int) (string, string) {
	date := time.Now().AddDate(0, 0, margin)
	day := date.Day()
	dayString := strconv.Itoa(day)

	month := int(date.Month())
	monthString := strconv.Itoa(month)
	if month < 10 {
		monthString = "0" + monthString
	}
	ex1 := dayString + "." + monthString
	if day < 10 {
		dayString = "0" + dayString
	}
	ex2 := dayString + "." + monthString
	if strings.Contains(data, ex1) || strings.Contains(data, ex2) {
		return ex1, ex2
	}
	return "", ""
}
func getSubgroupForDate(data, ex1, ex2 string) string {
	if strings.Contains(data, "/") {
		parts := strings.Split(data, "/")
		if len(parts) != 2 {
			return data
		}
		if strings.Contains(parts[0], ex1) || strings.Contains(parts[0], ex2) {
			return "[1 гр.]"
			//if isEven {
			//	return "[1 гр.]"
			//} else {
			//	return "[2 гр.]"
			//}
		} else if strings.Contains(parts[1], ex1) || strings.Contains(parts[1], ex2) {
			return "[2 гр.]"
			//if isEven {
			//	return "[2 гр.]"
			//} else {
			//	return "[1 гр.]"
			//}
		}
	} else {
		if strings.Contains(data, ex1) || strings.Contains(data, ex2) {
			return ex1
		}

	}
	return ""
}
func isContainsInDict(date string) bool {
	for _, v := range reservedDict {
		if v == date {
			return true
		}
	}

	return false
}

func (r *ScheduleRepository) GetCurrentDaySchedule(groupId int, margin int) ([]model.Lesson, time.Time) {
	day := time.Now().AddDate(0, 0, margin)
	dayNum := day.Weekday()

	groupSchedule, _ := r.GetScheduleByGroup(groupId)
	lessons := make([]model.Lesson, 0, 4)
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
	case dayNum == 0:
		lessons = []model.Lesson{}
	}

	return formScheduleList(lessons, margin), day
}

func GetScheduleStruct(body []byte) model.Schedule {
	var shed model.Schedule
	err := json.Unmarshal(body, &shed)
	if err != nil {
		log.Println(err)
	}
	return shed
}

func (r *ScheduleRepository) GetScheduleByGroup(group int) (model.Schedule, error) {
	var result string
	if err := r.store.db.QueryRow(
		"SELECT shedule FROM saved_timetable WHERE groupp = $1",
		group,
	).Scan(&result); err != nil {
		if err == sql.ErrNoRows {
			return model.Schedule{}, store.ErrRecordNotFound
		}

		return model.Schedule{}, err
	}
	scheduleStruct := GetScheduleStruct([]byte(result))
	return scheduleStruct, nil
}

func (r *ScheduleRepository) GetTeacherListStruct(groupId int) []model.Prepod {
	sched, _ := r.GetScheduleByGroup(groupId)
	prepodList := make([]model.Prepod, 0)
	v := reflect.ValueOf(sched)
	for i := 0; i < v.NumField(); i++ { // перебираем все поля структуры
		field := v.Field(i)
		if field.Kind() == reflect.Slice { // проверяем, что поле является срезом
			sliceValue := field.Interface()                   // получаем значение среза
			if slice, ok := sliceValue.([]model.Lesson); ok { // проверяем, что значение среза имеет тип []Lesson
				for _, lesson := range slice {
					lesson.PrepodName = strings.TrimSpace(lesson.PrepodName)
					lesson.DisciplType = strings.TrimSpace(lesson.DisciplType)
					lesson.DisciplName = strings.TrimSpace(lesson.DisciplName)
					added := false
					for k, prepod := range prepodList {
						if prepod.Name == lesson.PrepodName && prepod.Lesson == lesson.DisciplName {
							if !database.CheckInSlice(prepod.LessonType, lesson.DisciplType) {
								prepod.LessonType = append(prepod.LessonType, lesson.DisciplType)
								prepodList[k] = prepod
							}
							added = true
							break

						}
					}
					if added {
						continue
					}
					newPrepod := model.Prepod{
						LessonType: make([]string, 1),
						Name:       lesson.PrepodName,
						Lesson:     lesson.DisciplName,
					}
					newPrepod.LessonType[0] = lesson.DisciplType
					prepodList = append(prepodList, newPrepod)
				}
			}
		}
	}

	return prepodList

}
