package entry

import (
	"database/sql"
	"time"
)

type Entry struct {
	ID          int64
	UserID      int64
	Rate        int64
	Description sql.NullString
	CreatedAt   time.Time
}

func (e *Entry) Timestamp() string {
	return e.CreatedAt.Local().Format("Mon Jan 2 3:04PM")
}
