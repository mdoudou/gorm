package callbacks

import (
	"gorm.io/gorm"
)

func RowQuery(db *gorm.DB) {
	if db.Error == nil {
		if db.Statement.SQL.String() == "" {
			BuildQuerySQL(db)
		}

		if _, ok := db.Get("rows"); ok {
			db.Statement.Dest, db.Error = db.Statement.ConnPool.QueryContext(db.Statement.Context, db.Statement.SQL.String(), db.Statement.Vars...)
		} else {
			db.Statement.Dest = db.Statement.ConnPool.QueryRowContext(db.Statement.Context, db.Statement.SQL.String(), db.Statement.Vars...)
		}
	}
}
