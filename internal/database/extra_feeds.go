package database

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Norrun/feedmixer/internal/datautils"
	"github.com/Norrun/feedmixer/internal/display"
)

const getFeedAdv = `SELECT * FROM feeds
WHERE %s;`

const and = ", "
const or = " OR "
const include = " id IN "
const exclude = " id NOT IN "

const andBranchStart = `(SELECT feed_id
FROM tags_feeds
WHERE tag_id IN (`
const andBranchEnd = `)
GROUP BY feed_id
HAVING COUNT(tag_id) = %d)`

func (receiver *Queries) GetFeedsByTagTree(ctx context.Context, tree []display.Tag) ([]Feed, error) {
	var builder strings.Builder
	//correcting checkboxes, may change to do that earlier. Additionally need more advanced logic later to improve UX.
	tree = ProcessTagCheckTree(tree)

	// buildidng the query

	protree := tagTreeToAndPath(tree)

	for i, branch := range protree {
		if branch.V0 {
			builder.WriteString(include)
		} else {
			builder.WriteString(exclude)
		}
		builder.WriteString(andBranchStart)
		for j, leaf := range branch.V1 {
			builder.WriteString(strconv.Itoa(leaf))
			if j < len(branch.V1)-1 {
				builder.WriteString(and)
			}
		}
		builder.WriteString(fmt.Sprintf(andBranchEnd, len(branch.V1)))
		if i < len(protree)-1 {
			builder.WriteString(or)
		}
	}
	query := fmt.Sprintf(getFeedAdv, builder.String())
	fmt.Println(query)
	// Calling the database
	rows, err := receiver.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Feed
	for rows.Next() {
		var i Feed
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Url,
			&i.LastFetchedAt,
			&i.LastCheckedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil

}

func ProcessTagCheckTree(tags []display.Tag) []display.Tag {
	for i := range tags {
		tags[i] = ProcessTagCheck(tags[i])
	}
	return tags
}

func ProcessTagCheck(tag display.Tag) display.Tag {
	checked := tag.State

	for i := range tag.Related {
		subtag := ProcessTagCheck(tag.Related[i])
		if checked == display.Undecided && checked != subtag.State {
			if subtag.State == display.Mixed {
				checked = display.Mixed
			} else if checked == display.Excluded && subtag.State == display.Included {
				checked = display.Mixed
			}
		}
		tag.Related[i] = subtag

	}
	tag.State = checked
	return tag
}

func tagTreeToAndPath(tree []display.Tag) []datautils.Duo[bool, []int] {
	res := make([]datautils.Duo[bool, []int], 0, len(tree))
	for _, v := range tree {
		res = preProcessTagTreeRecur(v, make([]int, 0), res)
	}
	return res
}

func preProcessTagTreeRecur(node display.Tag, path []int, paths []datautils.Duo[bool, []int]) []datautils.Duo[bool, []int] {
	if node.State == display.Undecided {
		return paths
	}

	path = append(path, node.Id)
	if len(node.Related) == 0 {
		if node.State == display.Mixed {
			panic("Invalid tag forest, mixed leaf")
		}
		return append(paths, datautils.NewDuo(node.State == display.Included, path))
	}
	for _, child := range node.Related {
		paths = preProcessTagTreeRecur(child, path, paths)
	}
	return paths
}
