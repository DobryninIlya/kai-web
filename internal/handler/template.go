package handler

import (
	"bufio"
	"fmt"
	"log"
	service "main/internal/database"
	"os"
	"path/filepath"
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

// func GetMainView(data []interface{}) []byte {
func GetMainView(lessons []service.Lesson) []byte {
	tmp, _ := GetMainTemplate()
	allLessonData := ""
	for _, lesson := range lessons {
		prepodName := service.GetShortenName(lesson.PrepodName)
		allLessonData += getLesson([]interface{}{lesson.DayTime, "11:10", lesson.DisciplName, prepodName})

	}
	result := fmt.Sprintf(tmp, allLessonData)
	return []byte(result)
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
