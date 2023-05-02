package handler

import (
	"bufio"
	"fmt"
	"log"
	service "main/internal/database"
	"os"
	"path/filepath"
	"strings"
	"time"
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

func getLesson(data []interface{}) string {
	tmp, _ := GetLessonTemplate()
	result := fmt.Sprintf(tmp, data...)
	return result
}

// func GetMainView(date []interface{}) []byte {
func GetMainView(lessons []service.Lesson, date time.Time) []byte {
	tmp, _ := GetMainTemplate()

	tmpFilled := fmt.Sprintf(tmp, date.Format("Monday, 02 January"), "%v")
	if len(lessons) == 0 {
		return []byte(fmt.Sprintf(tmpFilled, service.GetNullDaySchedule()))
	}
	result := fmt.Sprintf(tmpFilled, GetLessonList(lessons))
	return []byte(result)
}

func GetLessonList(lessons []service.Lesson) string {
	allLessonData := ""
	nameStyle := ""
	for _, lesson := range lessons {
		//prepodName := service.GetShortenName(lesson.PrepodName)
		lessonDate := strings.TrimSpace(lesson.DayDate)
		room := service.GetRoom(lesson.AudNum)
		//lessonName := service.GetLessonName(lesson.DisciplName)
		lessonName := strings.TrimSpace(lesson.DisciplName)
		if len(lesson.DisciplName)/2 >= 30 {
			nameStyle = "font-size: 18px;"
		}
		allLessonData += getLesson([]interface{}{lesson.DayTime, room, nameStyle, lessonName, lessonDate})
		nameStyle = ""

	}
	//result := fmt.Sprintf(template, allLessonData)
	return allLessonData
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
