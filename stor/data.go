package stor

import "time"

type user struct {
	id_user int
}

type segment struct {
	id_segment int
	title      string
}

type subscription struct {
	id_subscription int
	id_user         int
	id_segment      int
	CreatedAt       time.Time
	ExpiresAt       time.Time
}
