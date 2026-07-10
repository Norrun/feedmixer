package database

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/Norrun/feedmixer/internal/datautils"
	"github.com/Norrun/feedmixer/internal/display"
)

const getTagsRelatedTags = `SELECT *
FROM tags t
WHERE t.id IN (
    SELECT a.tag_id
    FROM tags_feeds a
    WHERE a.feed_id IN (
        SELECT b.feed_id
        FROM tags_feeds b
        WHERE b.tag_id IN (%s)
        GROUP BY b.feed_id
        HAVING COUNT(DISTINCT b.tag_id) = ?
    )
)
AND t.id NOT IN (%s)`

func (receiver *Queries) GetTagsRelatedTags(ctx context.Context, tag_ids []int) ([]Tag, error) {
	tag_ids_str := strings.Repeat("?,", len(tag_ids)-1)
	tag_ids_str += "?"
	query := fmt.Sprintf(getTagsRelatedTags, tag_ids_str, tag_ids_str)
	values := datautils.AnySlice(tag_ids)
	values = append(values, len(tag_ids))
	values = slices.Concat(values, datautils.AnySlice(tag_ids))

	rows, err := receiver.db.QueryContext(
		ctx,
		query,
		values...,
	)
	if err != nil {
		return nil, err
	}
	var tags []Tag
	for rows.Next() {
		var tag Tag
		err = rows.Scan(
			&tag.ID,
			&tag.CreatedAt,
			&tag.UpdatedAt,
			&tag.Name,
			&tag.LastCheckedAt,
		)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (receiver *Queries) GetTagForest(ctx context.Context) ([]display.Tag, error) {
	dbtags, err := receiver.GetAllTags(ctx)
	if err != nil {
		return nil, err
	}
	/*
		tagsNfeeds, err := receiver.GetTagAndFeedIds(ctx)
		if err != nil {
			return nil, err
		}
		tagToFeeds, FeedToTags := make(map[int64][]int64, 0), make(map[int64][]int64, 0)

		for _, v := range tagsNfeeds {
			feeds := tagToFeeds[v.TagID]
			feeds = append(feeds, v.FeedID)
			tagToFeeds[v.TagID] = feeds

			tags := FeedToTags[v.FeedID]
			tags = append(tags, v.TagID)
			FeedToTags[v.FeedID] = tags
		}*/
	tags := make([]display.Tag, 0, len(dbtags))
	for _, v := range dbtags {
		tag, err := assembleTagTree(ctx, receiver, v, dbtags, []int64{v.ID})
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func assembleTagTree(ctx context.Context, q *Queries, tag Tag, tags []Tag, path []int64) (display.Tag, error) {

	relateddb, err := q.GetTagsRelatedTags(ctx, datautils.ConvertSlice(path, func(v int64) int { return int(v) }))
	if err != nil {
		return display.Tag{}, err
	}
	related, err := datautils.ConvertSliceErr(relateddb, func(t Tag) (display.Tag, error) { return assembleTagTree(ctx, q, t, tags, append(path, t.ID)) })
	if err != nil {
		return display.Tag{}, err
	}
	return display.Tag{Text: tag.Name, Id: int(tag.ID), Related: related}, nil

}

/*
func (receiver *Queries) GetAssembledTagTree(ctx context.Context) ([]display.Tag, error) {
	panic("missing dependency")
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

func (q *Queries) GetRelatedTags2(ctx context.Context) (any, any) {

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
}*/
