package database

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Norrun/feedmixer/internal/display"
)

func (receiver *Queries) GetAssembledTagTree(ctx context.Context) ([]display.Tag, error) {
	rows, err := receiver.GetTagTree(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error getting Row: %v", err)
	}

	result, err := receiver.newMethod(ctx, rows, 0)
	if err != nil {
		return nil, fmt.Errorf("Error Processing tag-tree: %v", err)
	}
	return result, nil
}

func (q *Queries) newMethod(ctx context.Context, rows []GetTagTreeRow, layer int) ([]display.Tag, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("Stopped: %v", ctx.Err())
	default:
		break
	}

	var tags []display.Tag

	for i, v := range rows {

		if v.Level < int64(layer) {
			return tags, nil
		}
		if v.Level > int64(layer) {
			continue
		}

		related, err := q.newMethod(ctx, rows[i:], layer+1)

		if err != nil {
			return nil, fmt.Errorf("Error with tag(%s) layer(%d): %v ", v.Name, layer, err)
		}

		tag := display.Tag{
			Text:    v.Name,
			Id:      strconv.Itoa(int(v.ID)),
			Checked: false,
			Related: related,
		}
		tags = append(tags, tag)

	}
	return nil, nil
}
