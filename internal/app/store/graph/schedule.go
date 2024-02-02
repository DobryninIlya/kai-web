package graph

import (
	"github.com/jmcvetta/neoism"
	"main/internal/app/model"
)

type Schedule struct {
	store *Store
}

type ScheduleInterface interface {
	AddSchedule(schedule model.Schedule)
}

func (s *Schedule) AddSchedule(schedule model.Schedule) {
	for _, lessons := range [][]model.Lesson{schedule.Day1, schedule.Day2, schedule.Day3, schedule.Day4, schedule.Day5, schedule.Day6} {
		for _, lesson := range lessons {
			// Добавляем каждый урок в базу данных
			if err := s.addLessonNode(lesson); err != nil {
				panic(err)
			}
		}
	}
}

func (s *Schedule) addLessonNode(lesson model.Lesson) error {
	// Создаем узлы для каждого поля урока
	nodes := make(map[string]string)
	nodes["dayDate"] = lesson.DayDate
	nodes["audNum"] = lesson.AudNum
	nodes["disciplName"] = lesson.DisciplName
	nodes["buildNum"] = lesson.BuildNum
	nodes["orgUnitName"] = lesson.OrgUnitName
	nodes["dayTime"] = lesson.DayTime
	nodes["dayNum"] = lesson.DayNum
	nodes["potok"] = lesson.Potok
	nodes["prepodName"] = lesson.PrepodName
	nodes["disciplNum"] = lesson.DisciplNum
	nodes["orgUnitId"] = lesson.OrgUnitId
	nodes["prepodLogin"] = lesson.PrepodLogin
	nodes["disciplType"] = lesson.DisciplType

	// Добавляем узлы, если их еще нет в базе данных
	for name, value := range nodes {
		if err := s.addLessonFieldNode(name, value); err != nil {
			return err
		}
	}

	// Создаем связи между узлами урока и его полями
	for name, _ := range nodes {
		if err := s.createRelation("Lesson", "LessonField", "HAS_"+name, lesson.DayDate); err != nil {
			return err
		}
	}

	return nil
}

func (s *Schedule) addLessonFieldNode(fieldName, fieldValue string) error {
	// Проверяем, существует ли уже узел для данного поля
	res := []struct {
		Field string `json:"f.field"`
	}{}
	cq := neoism.CypherQuery{
		Statement: `
            MATCH (f:LessonField {field: {field}})
            RETURN f.field
        `,
		Parameters: neoism.Props{"field": fieldValue},
		Result:     &res,
	}
	if err := s.store.db.Cypher(&cq); err != nil {
		return err
	}

	// Если узел не существует, создаем его
	if len(res) == 0 {
		cq := neoism.CypherQuery{
			Statement: `
                CREATE (f:LessonField {field: {field}})
            `,
			Parameters: neoism.Props{"field": fieldValue},
		}
		if err := s.store.db.Cypher(&cq); err != nil {
			return err
		}
	}

	return nil
}

func (s *Schedule) createRelation(startNodeLabel, endNodeLabel, relationType, dayDate string) error {
	// Создаем связь между узлами
	cq := neoism.CypherQuery{
		Statement: `
            MATCH (start:` + startNodeLabel + `), (end:` + endNodeLabel + `)
            WHERE start.dayDate = {dayDate}
            CREATE (start)-[:` + relationType + `]->(end)
        `,
		Parameters: neoism.Props{
			"dayDate": dayDate,
			//"fieldName": fieldName,
		},
	}
	if err := s.store.db.Cypher(&cq); err != nil {
		return err
	}
	return nil
}
