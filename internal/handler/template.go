package handler

import (
	"bufio"
	"fmt"
	"log"
	service "main/internal/database"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

const path = "templates"

func GetMainTemplate() (string, error) {
	data, err := readFile(filepath.Join("internal", path, "main.html"))
	if err != nil {
		log.Println(err)
		return "", err
	}
	return strings.Join(data, "\n"), nil
}
func GetRegistrationTemplate() (string, error) {
	data, err := readFile(filepath.Join("internal", path, "registration.html"))
	if err != nil {
		log.Println(err)
		return "", err
	}
	return strings.Join(data, "\n"), nil
}

func GetMainStylesheet() ([]byte, error) {
	data, err := readFile(filepath.Join("internal", path, "css", "main.css"))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return []byte(strings.Join(data, "\n")), nil
}

func GetLessonTemplate() (string, error) {
	data, err := readFile(filepath.Join("internal", path, "lesson.html"))
	if err != nil {
		log.Println(err)
		return "", err
	}
	//return []byte(strings.Join(data, "\n")), nil
	return strings.Join(data, "\n"), nil
}

func GetExamMainTemplate() (string, error) {
	data, err := readFile(filepath.Join("internal", path, "exam_main.html"))
	if err != nil {
		log.Println(err)
		return "", err
	}
	return strings.Join(data, "\n"), nil
}

func GetExamTemplate() (string, error) {
	data, err := readFile(filepath.Join("internal", path, "exam.html"))
	if err != nil {
		log.Println(err)
		return "", err
	}
	return strings.Join(data, "\n"), nil
}

func getLesson(data []interface{}) string {
	tmp, _ := GetLessonTemplate()
	result := fmt.Sprintf(tmp, data...)
	return result
}

// func GetMainView(date []interface{}) []byte {
// func GetMainView(lessons []service.Lesson, date time.Time) []byte {
func GetMainView() []byte {
	tmp, _ := GetMainTemplate()
	//tmpFilled := fmt.Sprintf(tmp, date.Format("Monday, 02 January"), "%v")
	//if len(lessons) == 0 {
	//	return []byte(fmt.Sprintf(tmpFilled, service.GetNullDaySchedule()))
	//}
	//result := fmt.Sprintf(tmpFilled, GetLessonList(lessons))
	return []byte(tmp)

}

func GetRegistrationView() []byte {
	tmp, _ := GetRegistrationTemplate()
	//tmpFilled := fmt.Sprintf(tmp, date.Format("Monday, 02 January"), "%v")
	//if len(lessons) == 0 {
	//	return []byte(fmt.Sprintf(tmpFilled, service.GetNullDaySchedule()))
	//}
	//result := fmt.Sprintf(tmpFilled, GetLessonList(lessons))
	return []byte(tmp)

}

func GetLessonList(lessons []service.Lesson) string {
	allLessonData := ""
	nameStyle := ""
	for _, lesson := range lessons {
		lessonTypeDiv := "<p class=\"lesson_type\" style=\"background-color: #%v\">%v</p>"
		//prepodName := service.GetShortenName(lesson.PrepodName)
		lessonDate := strings.TrimSpace(lesson.DayDate)
		room := service.GetRoom(lesson.AudNum)
		//lessonName := service.GetLessonName(lesson.DisciplName)
		lessonName := strings.TrimSpace(lesson.DisciplName)
		if len(lesson.DisciplName)/2 >= 30 {
			nameStyle = "font-size: 18px;"
		}
		lessonType := strings.TrimSpace(lesson.DisciplType)
		if lessonType == "лек" {
			lessonTypeDiv = fmt.Sprintf(lessonTypeDiv, "6CC241", "Лекция")
		} else if lessonType == "пр" {
			lessonTypeDiv = fmt.Sprintf(lessonTypeDiv, "E77A3D", "Практика")
		} else if lessonType == "л.р." {
			lessonTypeDiv = fmt.Sprintf(lessonTypeDiv, "3DD2E7", "Практика")
		} else {
			lessonTypeDiv = fmt.Sprintf(lessonTypeDiv, "3DA0E7", lessonType)
		}
		allLessonData += getLesson([]interface{}{lesson.DayTime, room, lessonTypeDiv, nameStyle, lessonName, lessonDate})
		nameStyle = ""

	}
	if len(lessons) == 0 {
		return service.GetNullDaySchedule()
	}
	//result := fmt.Sprintf(template, allLessonData)
	return allLessonData
}

func GetTeacherList(prepodList []service.Prepod) string {

}

func GetExamList(examElems []service.ExamStruct) string {
	mainTemplate, err := GetExamMainTemplate()
	if examElems == nil || err != nil {
		return mainTemplate
	}
	examsAllString := ""
	examElementTemplate, _ := GetExamTemplate()
	for _, exam := range examElems {
		prepodName := service.GetShortenName(exam.PrepodName)
		examsAllString += fmt.Sprintf(examElementTemplate, exam.ExamDate, exam.ExamTime, exam.DisciplName, prepodName, exam.AudNum, exam.BuildNum) + "\n"
	}
	return fmt.Sprintf(mainTemplate, examsAllString)
}

func IsEmptyStruct(s interface{}) bool {
	v := reflect.ValueOf(s)
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if !f.IsZero() {
			return false
		}
	}
	return true
}

func readFile(path string) ([]string, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("file does not exist")
			return nil, err
		}
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var rows []string
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		rows = append(rows, sc.Text())
	}
	return rows, nil

}
