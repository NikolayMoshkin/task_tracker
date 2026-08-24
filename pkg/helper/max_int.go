package helper

import (
	"errors"
	"task_tracker/pkg/model"
)

func MaxKey(arr []*model.Task) (int, error) {
	if len(arr) == 0 {
		return 0, errors.New("Map is empty")
	}

	var max int

	for _, t := range arr {
		if t.Id > max {
			max = t.Id
		}
	}

	return max, nil
}
