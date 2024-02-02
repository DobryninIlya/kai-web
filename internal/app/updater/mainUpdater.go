package updater_app

import (
	"main/internal/app/store/graph"
	"main/internal/app/updater/get_request"
)

func UpdateSchedule(store *graph.Store, groups []get_request.GroupInfo) {
	for _, group := range groups {
		schedule := get_request.GetScheduleByGroup(group.Id)
		store.Schedule().AddSchedule(schedule)
	}
}
