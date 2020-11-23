package model

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type BaseModel struct {
	ID         uint       `json:"id" gorm:"primaryKey"`
	CreateTime int64      `json:"create_time" gorm:"autoCreateTime"`
	UpdateTime int64      `json:"update_time" gorm:"autoUpdateTime"`
	DeleteTime DeleteTime `json:"delete_time"`
}

type DeleteTime sql.NullInt64

func Paginate(PageNo int, PageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch {
		case PageSize > 100:
			PageSize = 100
		case PageSize <= 0:
			PageSize = 10
		}
		switch {
		case PageNo <= 0:
			PageNo = 1
		}
		offset := (PageNo - 1) * PageSize
		return db.Offset(offset).Limit(PageSize)
	}
}

// Scan implements the Scanner interface.
func (n *DeleteTime) Scan(value interface{}) error {
	return (*sql.NullInt64)(n).Scan(value)
}

// Value implements the driver Valuer interface.
func (n DeleteTime) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Int64, nil
}

func (DeleteTime) QueryClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{
		clause.Where{Exprs: []clause.Expression{
			clause.Eq{
				Column: clause.Column{Table: clause.CurrentTable, Name: f.DBName},
				Value:  nil,
			},
		}},
	}
}

//func (v DeleteTime) UnmarshalJSON(data []byte) error {
//	// Unmarshalling into a pointer will let us detect null
//	var x *int64
//	if err := json.Unmarshal(data, &x); err != nil {
//		return err
//	}
//	if x != nil {
//		v.Valid = true
//		v.Int64 = *x
//	} else {
//		v.Valid = false
//	}
//	return nil
//}

type SoftDeleteQueryClause struct {
	Field *schema.Field
}

func (sd SoftDeleteQueryClause) Name() string {
	return ""
}

func (sd SoftDeleteQueryClause) Build(clause.Builder) {
}

func (sd SoftDeleteQueryClause) MergeClause(*clause.Clause) {
}

func (sd SoftDeleteQueryClause) ModifyStatement(stmt *gorm.Statement) {
	if _, ok := stmt.Clauses["soft_delete_enabled"]; !ok {
		stmt.AddClause(clause.Where{Exprs: []clause.Expression{
			clause.Eq{Column: clause.Column{Table: clause.CurrentTable, Name: sd.Field.DBName}, Value: nil},
		}})
		stmt.Clauses["soft_delete_enabled"] = clause.Clause{}
	}
}

func (DeleteTime) DeleteClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{SoftDeleteDeleteClause{Field: f}}
}

type SoftDeleteDeleteClause struct {
	Field *schema.Field
}

func (sd SoftDeleteDeleteClause) Name() string {
	return ""
}

func (sd SoftDeleteDeleteClause) Build(clause.Builder) {
}

func (sd SoftDeleteDeleteClause) MergeClause(*clause.Clause) {
}

func (sd SoftDeleteDeleteClause) ModifyStatement(stmt *gorm.Statement) {
	if stmt.SQL.String() == "" {
		//stmt.AddClause(clause.Set{{Column: clause.Column{Name: sd.Field.DBName}, Value: stmt.DB.NowFunc()}})
		stmt.AddClause(clause.Set{{Column: clause.Column{Name: sd.Field.DBName}, Value: time.Now().Unix()}})

		if stmt.Schema != nil {
			_, queryValues := schema.GetIdentityFieldValuesMap(stmt.ReflectValue, stmt.Schema.PrimaryFields)
			column, values := schema.ToQueryValues(stmt.Table, stmt.Schema.PrimaryFieldDBNames, queryValues)

			if len(values) > 0 {
				stmt.AddClause(clause.Where{Exprs: []clause.Expression{clause.IN{Column: column, Values: values}}})
			}

			if stmt.ReflectValue.CanAddr() && stmt.Dest != stmt.Model && stmt.Model != nil {
				_, queryValues = schema.GetIdentityFieldValuesMap(reflect.ValueOf(stmt.Model), stmt.Schema.PrimaryFields)
				column, values = schema.ToQueryValues(stmt.Table, stmt.Schema.PrimaryFieldDBNames, queryValues)

				if len(values) > 0 {
					stmt.AddClause(clause.Where{Exprs: []clause.Expression{clause.IN{Column: column, Values: values}}})
				}
			}
			fmt.Println(values)
		}
		stmt.AddClauseIfNotExists(clause.Update{})
		stmt.Build("UPDATE", "SET", "WHERE")
	}
}
