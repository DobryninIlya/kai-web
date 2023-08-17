package tools

import (
	"bufio"
	"fmt"
	"log"
	"main/internal/app/database"
	"main/internal/app/model"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

const path = "templates"

func GetMainTemplate() (string, error) {
	data, err := readFile(filepath.Join("internal", "app", path, "main.html"))
	if err != nil {
		log.Println(err)
		return "", err
	}
	return strings.Join(data, "\n"), nil
}
func GetRegistrationTemplate() (string, error) {
	data, err := readFile(filepath.Join("internal", "app", path, "registration.html"))
	if err != nil {
		log.Println(err)
		return "", err
	}
	return strings.Join(data, "\n"), nil
}

func GetRegistrationIDcardTemplate() (string, error) {
	data, err := readFile(filepath.Join("internal", "app", path, "registrationIDcard.html"))
	if err != nil {
		log.Println(err)
		return "", err
	}
	return strings.Join(data, "\n"), nil
}

func GetMainStylesheet() ([]byte, error) {
	data, err := readFile(filepath.Join("internal", "app", path, "css", "main.css"))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return []byte(strings.Join(data, "\n")), nil
}

func GetLessonTemplate() (string, error) {
	data, err := readFile(filepath.Join("internal", "app", path, "lesson.html"))
	if err != nil {
		log.Println(err)
		return "", err
	}
	//return []byte(strings.Join(data, "\n")), nil
	return strings.Join(data, "\n"), nil
}

func GetExamMainTemplate() (string, error) {
	data, err := readFile(filepath.Join("internal", "app", path, "exam_main.html"))
	if err != nil {
		log.Println(err)
		return "", err
	}
	return strings.Join(data, "\n"), nil
}

func GetExamTemplate() (string, error) {
	data, err := readFile(filepath.Join("internal", "app", path, "exam.html"))
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

func GetMainView() []byte {
	tmp, _ := GetMainTemplate()
	return []byte(tmp)

}

func GetRegistrationView() []byte {
	tmp, _ := GetRegistrationTemplate()
	return []byte(tmp)
}

func GetLessonList(lessons []model.Lesson) string {
	allLessonData := ""
	nameStyle := ""
	for _, lesson := range lessons {
		lessonTypeDiv := "<p class=\"lesson_type\" style=\"background-color: #%v\">%v</p>"
		//prepodName := service.GetShortenName(lesson.PrepodName)
		lessonDate := strings.TrimSpace(lesson.DayDate)
		room := database.GetRoom(lesson.AudNum)
		lesson.BuildNum = strings.TrimSpace(lesson.BuildNum)
		if len(lesson.BuildNum) < 3 {
			room = lesson.BuildNum + "зд. " + room
		}
		//lessonName := service.GetLessonName(lesson.DisciplName)
		lessonName := strings.TrimSpace(lesson.DisciplName)
		//if len(lesson.DisciplName)/2 >= 30 {
		//	nameStyle = "font-size: 18px;"
		//}
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
		allLessonData += getLesson([]interface{}{
			lesson.DayTime, room, lessonTypeDiv, nameStyle, database.GetShortenLessonName(lessonName), lessonDate,
			database.GetShortenLessonName(lessonName), lessonName, lesson.PrepodName, fmt.Sprintf("%v здание, %v ауд.", lesson.BuildNum, lesson.AudNum), lessonTypeDiv,
		})
		nameStyle = ""

	}
	if len(lessons) == 0 {
		return database.GetNullDaySchedule()
	}
	//result := fmt.Sprintf(template, allLessonData)
	return allLessonData
}

func GetMainTeachersTemplate() (string, error) {
	data, err := readFile(filepath.Join("internal", "app", path, "teachers_main.html"))
	if err != nil {
		log.Println(err)
		return "", err
	}
	return strings.Join(data, "\n"), nil
}
func GetTeachersTemplate() (string, error) {
	data, err := readFile(filepath.Join("internal", "app", path, "teachers.html"))
	if err != nil {
		log.Println(err)
		return "", err
	}
	return strings.Join(data, "\n"), nil
}

func GetTeacherList(prepodList []model.Prepod) string {
	mainTemplate, err := GetMainTeachersTemplate()
	if err != nil {
		log.Printf("Ошибка teachers template %v", err)
	}
	teacherDiv, err := GetTeachersTemplate()
	if err != nil {
		log.Printf("Ошибка teachers div %v", err)
	}
	result := ""
	for _, prepod := range prepodList {
		if prepod.Name == "" {
			continue
		}
		lessonTypes := ""
		for _, s := range prepod.LessonType {
			color := "#4542DE"
			switch s {
			case "пр":
				color = "#E77A3D"
			case "л.р.":
				color = "#3DD2E7"
			case "лек":
				color = "#6CC241"
			}
			lessonTypes += fmt.Sprintf("<p class=\"lesson_type\" style=\"background-color: %v;\">%v</p>", color, s)
		}
		result += fmt.Sprintf(teacherDiv, lessonTypes, prepod.Name, prepod.Lesson)
	}
	return fmt.Sprintf(mainTemplate, result)

}

func GetScoreMainTemplate() (string, error) {
	data, err := readFile(filepath.Join("internal", "app", path, "score_main.html"))
	if err != nil {
		log.Println(err)
		return "", err
	}
	return strings.Join(data, "\n"), nil
}

func GetScoreTemplate() (string, error) {
	data, err := readFile(filepath.Join("internal", "app", path, "score.html"))
	if err != nil {
		log.Println(err)
		return "", err
	}
	return strings.Join(data, "\n"), nil
}

func GetScoreList(scoreList []ScoreElement) string {
	scoreTemplate, err := GetScoreMainTemplate()
	if scoreList == nil || err != nil {
		return scoreTemplate
	}
	scoreAllString := ""
	scoreElementTemplate, _ := GetScoreTemplate()
	for _, elem := range scoreList {
		scoreAllString += fmt.Sprintf(scoreElementTemplate, elem.Name, elem.PreviouslyScore) + "\n"
	}
	return fmt.Sprintf(scoreTemplate, scoreAllString)
}

func GetExamList(examElems []model.Exam) string {
	mainTemplate, err := GetExamMainTemplate()
	if examElems == nil || err != nil {
		return mainTemplate
	}
	examsAllString := ""
	examElementTemplate, _ := GetExamTemplate()
	if len(examElems) == 0 {
		examsAllString += fmt.Sprintf(examElementTemplate, "", "", "Данные не найдены", "", "", "")
		return fmt.Sprintf(mainTemplate, examsAllString)
	}
	for _, exam := range examElems {
		prepodName := database.GetShortenName(exam.PrepodName)
		examsAllString += fmt.Sprintf(examElementTemplate, exam.ExamDate, exam.ExamTime, exam.DisciplName, prepodName, exam.AudNum, exam.BuildNum) + "\n"
	}
	return fmt.Sprintf(mainTemplate, examsAllString)
}

func GetRegistrationIDcard() []byte {
	tmp, _ := GetRegistrationIDcardTemplate()
	return []byte(tmp)
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
