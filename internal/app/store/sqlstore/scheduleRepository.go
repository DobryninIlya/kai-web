package sqlstore

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"main/internal/app/formatter"
	"main/internal/app/model"
	"main/internal/app/store"
	_ "main/internal/app/store"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var reservedDict = []string{"чет", "неч", "чет/неч", "неч/чет"}

// ScheduleRepository ...
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
	).Scan(&result); err != nil || len(result) == 0 {
		if err == sql.ErrNoRows || len(result) == 0 {
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
	return 0, errors.New("Расписание не найдено для группы: " + strconv.Itoa(id))
}

func (r ScheduleRepository) MarkDeletedLesson(creator string, groupid int, lessonId int, uniqString string, platform string) (int, error) {
	var result int
	if err := r.store.db.QueryRow(
		"INSERT INTO public.deleted_lessons(groupid, creator, creator_platform, lesson_id, date, uniqstring) VALUES ($1, $2, $3, $4, NOW(), $5) RETURNING id;",
		groupid,
		creator,
		platform,
		lessonId,
		uniqString,
	).Scan(&result); err != nil {
		return 0, err
	}
	return result, nil
}

func (r ScheduleRepository) ReturnDeletedLesson(lessonId int, uniqString string) (int, error) {
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

func formScheduleList(lessons []model.Lesson, margin int, weekParity int) []model.Lesson {
	_, week := time.Now().AddDate(0, 0, margin).ISOWeek()
	result := make([]model.Lesson, 0)
	isEven := (week%2 + weekParity) == 0 // Булевый параметр четности недели.
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
		if date == "чет" && isEven { // Если дата чет и неделя четная - добавляем в список
			result = append(result, lesson)
		} else if date == "неч" && !isEven { // Если дата неч и неделя нечетная - добавляем в список
			result = append(result, lesson)
		} else if (date == "чет/неч" && isEven) || (date == "неч/чет" && !isEven) { // Если дата чет/неч и неделя четная - добавляем в список 1 группу
			lesson.DayDate = "[1 гр.]" // Если дата неч/чет и неделя нечетная - добавляем в список 1 группу
			result = append(result, lesson)
		} else if (date == "неч/чет" && isEven) || (date == "чет/неч" && !isEven) { // Если дата неч/чет и неделя четная - добавляем в список 2 группу
			lesson.DayDate = "[2 гр.]" // Если дата чет/неч и неделя нечетная - добавляем в список 2 группу
			result = append(result, lesson)
		} else if ex1, ex2 := isContainDate(date, margin); ex1+ex2 != "" { // Проверяем содержится ли дата в строке и хотя бы одна из дат совпадает с текущей (включая отступ)
			lesson.DayDate = getSubgroupForDate(date, ex1, ex2) // Если содержится - добавляем в список необходимую подгруппу
			result = append(result, lesson)
		} else if !isContainsInDict(date) { // Проверяем содержится ли дата в словаре зарезервированных слов (чет, нечет, чет/неч, неч/чет)
			result = append(result, lesson) // Если не содержится - добавляем в список
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

func (r ScheduleRepository) GetCurrentDaySchedule(groupId int, margin int, weekParity int) ([]model.Lesson, time.Time, error) {
	day := time.Now().AddDate(0, 0, margin)
	dayNum := day.Weekday()

	groupSchedule, err := r.GetScheduleByGroup(groupId)
	if err != nil {
		return nil, time.Time{}, err
	}
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

	return formScheduleList(lessons, margin, weekParity), day, nil
}

func GetScheduleStruct(body []byte) model.Schedule {
	re := regexp.MustCompile("\\s+")
	//re := regexp.MustCompile("\\s{2,}\"")
	//re := regexp.MustCompile("\\s+(?=\").|(?<=\\s)\"")
	s := string(bytes.TrimSpace(body))
	//s := string(bytes.TrimRight(body, " "))
	s = re.ReplaceAllString(s, ` `)
	body = []byte(s)
	var shed model.Schedule
	err := json.Unmarshal(body, &shed)
	if err != nil {
		log.Println(err)
	}
	val := reflect.ValueOf(&shed).Elem()
	trimRightSpaces(val)
	return shed
}

func trimRightSpaces(value reflect.Value) {
	switch value.Kind() {
	case reflect.Ptr:
		trimRightSpaces(value.Elem())
	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			field := value.Field(i)
			trimRightSpaces(field)
		}
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			trimRightSpaces(value.Index(i))
		}
	case reflect.String:
		if value.CanSet() {
			newValue := strings.TrimRight(value.String(), " ")
			value.SetString(newValue)
		}
	}
}

func (r ScheduleRepository) GetScheduleByGroup(group int) (model.Schedule, error) {
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

func (r ScheduleRepository) GetTeacherListStruct(groupId int) ([]model.Prepod, error) {
	sched, err := r.GetScheduleByGroup(groupId)
	if err != nil {
		return nil, err
	}
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
							if !formatter.CheckInSlice(prepod.LessonType, lesson.DisciplType) {
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

	return prepodList, nil
}

func (r ScheduleRepository) GetScheduleWithDeletedLessons(schedule model.Schedule, groupId int) (model.Schedule, error) {
	deletedLessons, err := r.GetDeletedLessonsByGroup(groupId)
	if err != nil {
		return model.Schedule{}, err
	}
	day1, _ := checkDeletedLessonsInDay(schedule.Day1, deletedLessons)
	schedule.Day1 = day1
	day2, _ := checkDeletedLessonsInDay(schedule.Day2, deletedLessons)
	schedule.Day2 = day2
	day3, _ := checkDeletedLessonsInDay(schedule.Day3, deletedLessons)
	schedule.Day3 = day3
	day4, _ := checkDeletedLessonsInDay(schedule.Day4, deletedLessons)
	schedule.Day4 = day4
	day5, _ := checkDeletedLessonsInDay(schedule.Day5, deletedLessons)
	schedule.Day5 = day5
	day6, _ := checkDeletedLessonsInDay(schedule.Day6, deletedLessons)
	schedule.Day6 = day6

	return schedule, nil
}

func checkDeletedLessonsInDay(lessons []model.Lesson, deletedLessons []model.DeletedLessonsMin) ([]model.Lesson, error) {
	for i, lesson := range lessons {
		lessonType := strings.TrimSpace(lesson.DisciplType)
		uniqstring := lessonType + "_" + strings.TrimSpace(lesson.DayTime) + "_" + strings.TrimSpace(lesson.DayDate)
		disciplNum, err := strconv.Atoi(lesson.DisciplNum)
		if err != nil {
			return nil, err
		}
		for _, deletedLesson := range deletedLessons {
			if deletedLesson.LessonId == disciplNum && strings.TrimSpace(deletedLesson.Uniqstring) == uniqstring {
				lessons[i].MarkedDeleted = true
			}
		}
	}
	return lessons, nil
}
