package database

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"regexp"

	"github.com/jmoiron/sqlx"
)

type QueryBuilder struct {
	BaseSQL         string // initial sql string to build query from
	Subquery        bool   // wraps query in parenthesis
	Single          bool   // returns single entity
	Paginate        bool
	ItemsPerPage    int
	PageNum         int // index from 1
	OrderBy         []string
	OrderDescending bool
	GroupBy         []string
	Where           []QueryFilter
	Setters         []QuerySetter // used for insert and update compilation
}

type QueryFilter struct {
	Unsafe          bool // disables validation if true
	Column          string
	Operator        int           // const enum
	Parameter       interface{}   // not used if subquery not nil
	SubqueryBuilder *QueryBuilder // builds a SELECT subquery
}

type QueryBetween struct {
	First  interface{}
	Second interface{}
}

type QuerySetter struct {
	Column          string        // column to set
	Parameter       interface{}   // parameter to set IF THERE IS NO SUBQUERY
	SubqueryBuilder *QueryBuilder // builds a SELECT subquery - NOTE - if the subquery contains filters, those filters will automatically append to the parameter list for the generated query
}

// filter operators
const (
	EQ      = iota // equal
	NE      = iota // not equal
	GT      = iota // greater than
	LT      = iota // less than
	GE      = iota // greater than or equal to
	LE      = iota // less than or equal to
	LIKE    = iota // like
	BETWEEN = iota // between
)

func Update[T any](qb *QueryBuilder, db *sqlx.DB, manualParams ...interface{}) (sql.Result, error) {
	sql, params, buildErr := buildUpdate[T](qb)
	if buildErr != nil {
		return nil, buildErr
	}

	params = append(params, manualParams...)

	return db.Exec(sql, params...)
}

func Insert[T any](qb *QueryBuilder, db *sqlx.DB, manualParams ...interface{}) (sql.Result, error) {
	sql, params, buildErr := buildInsert[T](qb)
	if buildErr != nil {
		return nil, buildErr
	}
	params = append(params, manualParams...)

	return db.Exec(sql, params...)
}

func Delete[T any](qb *QueryBuilder, db *sqlx.DB, manualParams ...interface{}) (sql.Result, error) {
	sql, params, buildErr := buildDelete[T](qb)
	if buildErr != nil {
		return nil, buildErr
	}

	params = append(params, manualParams...)

	return db.Exec(sql, params...)
}

func Select[T any](qb *QueryBuilder, db *sqlx.DB, manualParams ...interface{}) ([]T, error) {
	var entities []T

	sql, params, buildErr := buildSelect[T](qb)
	if buildErr != nil {
		return nil, buildErr
	}

	params = append(params, manualParams...)

	err := db.Select(&entities, sql, params...)

	return entities, err
}

func Get[T any](qb *QueryBuilder, db *sqlx.DB, manualParams ...interface{}) (T, error) {
	var entity T

	qb.Single = true

	sql, params, buildErr := buildSelect[T](qb)
	if buildErr != nil {
		return entity, buildErr
	}

	params = append(params, manualParams...)

	err := db.Get(&entity, sql, params...)

	return entity, err
}

// internal functions

func buildWhere[T any](qb *QueryBuilder) (string, []interface{}, error) {
	var sql string

	sql += " "

	if qb == nil {
		return "", nil, errors.New("no query filter provided")
	}

	if len(qb.Where) == 0 {
		return "", nil, nil
	}

	var params []interface{}

	for i, v := range qb.Where {
		if !v.Unsafe && !validateSQLColName[T](v.Column) {
			return "", nil, errors.New("invalid column name provided in WHERE clause. column: " + v.Column)
		}

		if i == 0 {
			sql += "WHERE "
		} else {
			sql += "AND "
		}

		sql += v.Column + " "

		// single param operators
		switch v.Operator {
		case EQ:
			sql += "= "
		case NE:
			sql += "<> "
		case GT:
			sql += "> "
		case LT:
			sql += "< "
		case GE:
			sql += ">= "
		case LE:
			sql += "<= "
		case LIKE:
			sql += "LIKE "
		case BETWEEN:
			sql += "BETWEEN ? AND ? "
		}

		if v.SubqueryBuilder == nil && v.Operator != BETWEEN {
			sql += "? "
		}

		if v.Parameter == nil {
			return "", nil, errors.New("WHERE clause not nil, but parameters nil for column: " + v.Column)
		}

		between, isBetween := v.Parameter.(QueryBetween)

		if v.Operator == BETWEEN {
			// between does not work with subqueries atm
			if !isBetween {
				return "", nil, errors.New("Filter parameter must be of type `QueryBetween` when using the BETWEEN operator. Column: " + v.Column)
			}

			params = append(params, between.First)
			params = append(params, between.Second)
		} else {
			if isBetween {
				return "", nil, errors.New("Attempt to use a between filter for a non between query. Column: " + v.Column)
			}

			if v.SubqueryBuilder != nil {
				// recurse!
				rSql, rParams, subqueryErr := buildSelect[T](v.SubqueryBuilder)

				if subqueryErr != nil {
					return "", nil, subqueryErr
				}

				sql += rSql + " "
				params = append(params, rParams...)
			} else {
				params = append(params, v.Parameter)
			}
		}
	}

	return sql, params, nil
}

func buildSelect[T any](qb *QueryBuilder) (string, []interface{}, error) {
	if qb == nil {
		return "", nil, errors.New("no query filter provided")
	}

	sql := qb.BaseSQL

	var params []interface{}

	sql += " "

	// where filters
	wSql, wParams, wErr := buildWhere[T](qb)
	if wErr != nil {
		return "", nil, wErr
	}

	params = append(params, wParams...)
	sql += wSql

	// grouping
	if len(qb.GroupBy) > 0 {
		sql += "GROUP BY "

		for i, v := range qb.GroupBy {
			if !validateSQLColName[T](v) {
				return "", nil, errors.New("invalid name for groupby clause")
			}

			if i == (len(qb.GroupBy) - 1) {
				sql += v + " "
			} else {
				sql += v + ", "
			}
		}
	}

	// ordering
	if len(qb.OrderBy) > 0 {
		sql += "ORDER BY "

		for i, v := range qb.OrderBy {
			if !validateSQLColName[T](v) {
				return "", nil, errors.New("invalid name for orderby clause")
			}

			if i == (len(qb.OrderBy) - 1) {
				sql += v + " "
			} else {
				sql += v + ", "
			}
		}

		if qb.OrderDescending {
			sql += "DESC "
		} else {
			sql += "ASC "
		}
	}

	// return first result
	if qb.Single {
		sql += "LIMIT 1"
	}

	// pagination
	if !qb.Single && qb.Paginate {
		if qb.PageNum <= 0 {
			qb.PageNum = 1
		}

		if qb.ItemsPerPage <= 0 {
			qb.ItemsPerPage = 10
		}

		sql += fmt.Sprintf("LIMIT %d ", qb.ItemsPerPage)

		offset := (qb.PageNum - 1) * qb.ItemsPerPage
		sql += fmt.Sprintf("OFFSET %d", offset)
	}

	if qb.Subquery {
		sql = "(" + sql + ")"
	}

	return sql, params, nil
}

func buildInsert[T any](qb *QueryBuilder) (string, []interface{}, error) {
	sql := qb.BaseSQL

	sql += " ("

	var columns []string
	var parameters []interface{}

	for _, k := range qb.Setters {
		columns = append(columns, k.Column)
	}

	if len(columns) == 0 {
		return "", nil, errors.New("one or more setters contains invalid column name")
	}

	for i, v := range columns {
		if i == len(columns)-1 {
			sql += v
		} else {
			sql += v + ","
		}
	}

	sql += ") VALUES ("

	for i, v := range qb.Setters {
		var subSql string
		var subParams []interface{}
		var buildSubErr error

		if v.SubqueryBuilder != nil {
			subSql, subParams, buildSubErr = buildSelect[T](v.SubqueryBuilder)
			if buildSubErr != nil {
				return "", nil, buildSubErr
			}

			parameters = append(parameters, subParams...)
		} else {
			parameters = append(parameters, v.Parameter)
		}

		if i == len(columns)-1 {
			if v.SubqueryBuilder != nil {
				sql += subSql
			} else {
				sql += "?"
			}
		} else {
			if v.SubqueryBuilder != nil {
				sql += subSql + ","
			} else {
				sql += "?,"
			}
		}
	}

	sql += ") "

	return sql, parameters, nil
}

func buildUpdate[T any](qb *QueryBuilder) (string, []interface{}, error) {
	sql := qb.BaseSQL

	sql += " "

	var params []interface{}

	for i, v := range qb.Setters {
		var subSql string
		var subParams []interface{}
		var buildSubErr error

		if v.SubqueryBuilder != nil {
			subSql, subParams, buildSubErr = buildSelect[T](v.SubqueryBuilder)
			if buildSubErr != nil {
				return "", nil, buildSubErr
			}

			params = append(params, subParams...)
		} else {
			params = append(params, v.Parameter)
		}

		if i == len(qb.Setters)-1 && i == 0 {
			if v.SubqueryBuilder != nil {
				sql += "SET " + v.Column + " = " + subSql + " "
			} else {
				sql += "SET " + v.Column + " = " + "? "
			}
		} else if i == 0 {
			if v.SubqueryBuilder != nil {
				sql += "SET " + v.Column + " = " + subSql + ", "
			} else {
				sql += "SET " + v.Column + " = " + "?, "
			}
		} else if i == len(qb.Setters)-1 {
			if v.SubqueryBuilder != nil {
				sql += v.Column + " = " + subSql + " "
			} else {
				sql += v.Column + " = " + "? "
			}
		} else {
			if v.SubqueryBuilder != nil {
				sql += v.Column + " = " + subSql + ", "
			} else {
				sql += v.Column + " = " + "?, "
			}
		}
	}

	// where filters
	wSql, wParams, wErr := buildWhere[T](qb)
	if wErr != nil {
		return "", nil, wErr
	}

	params = append(params, wParams...)
	sql += wSql

	return sql, params, nil
}

func buildDelete[T any](qb *QueryBuilder) (string, []interface{}, error) {
	sql := qb.BaseSQL

	outSql, params, err := buildWhere[T](qb)
	if err != nil {
		return "", nil, err
	}

	outSql = sql + outSql

	return outSql, params, nil
}

func validateSQLColName[T any](input string) bool {
	// only match alphanumerics and underscores
	// will still match if dot '.' present, but only if
	// the dot is between two alphanumerics
	valid := regexp.MustCompile(`^[A-Za-z0-9_]+(\.[A-Za-z0-9_]+)*$`)
	if !valid.MatchString(input) {
		return false
	}

	// Check if the column exists as a `db` note for the given type T
	var t T
	typeOfT := reflect.TypeOf(t)

	for i := 0; i < typeOfT.NumField(); i++ {
		field := typeOfT.Field(i)
		if dbTag, ok := field.Tag.Lookup("db"); ok && dbTag == input {
			return true
		}
	}
	return false
}
