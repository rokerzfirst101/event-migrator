package interview

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"event-migration-script/store"
)

type interview struct {
}

func New() store.Interview {
	return &interview{}
}

func (i interview) GetEventIDMap(ctx *gofr.Context) (map[int]string, error) {
	res := make(map[int]string)

	rows, err := ctx.DB().QueryContext(ctx, getEventIDMapQuery)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		interviewID := 0
		eventID := ""

		err = rows.Scan(&interviewID, &eventID)
		if err != nil {
			return nil, err
		}

		res[interviewID] = eventID
	}

	return res, nil
}

func (i interview) UpdateEventID(ctx *gofr.Context, oldEventID, newEventID string) error {
	_, err := ctx.DB().ExecContext(ctx, updateEventIDQuery, newEventID, oldEventID)
	if err != nil {
		return err
	}

	return nil
}
