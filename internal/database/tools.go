package database

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const lessonName = 20

var roomsDictionary = make(map[string]string)

func init() {
	roomsDictionary["кск каи олимп"] = "Олимп"
	// тут будет много текста
}

func GetShortenName(name string) string {
	if len(name) < 3 { // Либо пустая, либо слишком короткая
		return ""
	}
	parts := strings.Split(strings.TrimSpace(name), " ")
	name1, _ := utf8.DecodeRuneInString(parts[1])
	surname1, _ := utf8.DecodeRuneInString(parts[2])
	return fmt.Sprintf("%v %v.%v.", parts[0], string(name1), string(surname1))
}

func cut(text string, limit int) string {
	runes := []rune(text)
	if len(runes) >= limit {
		return string(runes[:limit])
	}
	return text
}

func GetRoom(room string) string {
	room = strings.TrimSpace(room)
	res := roomsDictionary[strings.ToLower(room)]
	if res != "" {
		return res
	}
	if len(room) > 50 {
		room = cut(room, 15)
		return room
	}

	if len(room) > 5 {
		result := ""
		parts := strings.Split(room, " ")
		for _, part := range parts {
			result += cut(part, 3)
			if len(part) > 3 {
				result += ". "
			} else {
				result += " "
			}
		}
		return strings.TrimSpace(result)
	}
	return room
}

func GetLessonName(name string) string {
	name = strings.TrimSpace(name)
	res := roomsDictionary[strings.ToLower(name)]
	if res != "" {
		return res
	}
	if len(name) < 21 {
		return name
	}

	if len(name) < 100 { // Обрезаем каждое слово так, чтобы в целом все влезло
		result := ""
		parts := strings.Split(name, " ")
		cutDistance := lessonName / len(parts)
		for _, part := range parts {
			result += cut(part, cutDistance)
			if len(part) > 3 {
				result += ". "
			} else {
				result += " "
			}
		}
		return strings.TrimSpace(result)
	}
	return cut(name, 20)
}

func GetNullDaySchedule() string {
	return "<div class=\"lesson_list\" id=\"lesson_list\">\n    <div class=\"lesson_none\"> <p>Занятий не найдено</p></div>\n</div>"
}
