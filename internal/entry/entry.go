package entry

import "database/sql"

type Entry struct {
	ID          int64
	UserID      int64
	Rate        int64
	Description sql.NullString
}
