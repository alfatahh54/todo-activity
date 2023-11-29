package db

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/alfatahh54/todo-activity/utils"
	_ "github.com/go-sql-driver/mysql"
)

func DbConnect(dbUser, dbPass, dbHost, dbPort, dbName string) *sql.DB {
	dbAccess := dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?parseTime=true&loc=Asia%2FJakarta"
	db, err := sql.Open("mysql", dbAccess)
	if err != nil {
		panic(err.Error())
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db
}

var DB *sql.DB
var Db = &SqlDB{}

func Open() *sql.DB {
	dbUser := utils.GoDotEnvVariable("DB_USER")
	dbPass := utils.GoDotEnvVariable("DB_PASS")
	dbName := utils.GoDotEnvVariable("DB_NAME")
	dbHost := utils.GoDotEnvVariable("DB_HOST")
	dbPort := utils.GoDotEnvVariable("DB_PORT")
	dbDriver := utils.GoDotEnvVariable("DB_DRIVER")

	println(dbHost, dbUser, dbName, dbPass, dbPort)

	if dbPort == "" {
		dbPort = "3306"
	}

	if dbDriver == "" {
		dbDriver = "mysql"
	}

reconnect:
	if dbDriver == "mysql" {
		DB = DbConnect(dbUser, dbPass, dbHost, dbPort, dbName)
	}

	err := DB.Ping()
	if err != nil {
		fmt.Println("Error Pinging DB : ", err)
		// panic(err)
		goto reconnect
	}

	fmt.Println("Connected to db!")

	return DB
}

func init() {
	if utils.TestMode() {
		return
	}
	DB = Open()
	Db.DB = DB
}
func MysqlQuery(query string, name string, fields []string, params ...interface{}) []interface{} {
	q := query
	// start := time.Now()
	for {
		if DB != nil {
			break
		}
	}
	rows, err := DB.Query(q, params...)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	cols, _ := rows.Columns()

	// var getData [][]string
	var mysqlData []interface{}
	// getData = append(getData, fields)
	for rows.Next() {
		val := make([]sql.NullString, len(cols))
		values := make([]interface{}, len(cols))
		for i := range cols {
			values[i] = &val[i]
		}

		if err := rows.Scan(values...); err != nil {
			panic(err)
		}
		rowResult := make(map[string]interface{})

		columNames, _ := rows.Columns()

		for i, v := range val {
			if columNames[i] == "1" {
				continue
			}
			rowResult[columNames[i]] = v.String
		}
		mysqlData = append(mysqlData, rowResult)
	}
	return mysqlData
}

func MysqlQuerySingleRow(query string, name string, fields any, params ...interface{}) error {
	q := query
	// start := time.Now()
	rows, err := DB.Query(q, params...)
	if err != nil {
		fmt.Println("Error query SingleRow : ", err)
		return err
	}

	defer rows.Close()

	cols, _ := rows.Columns()

	rowResult := make(map[string]interface{})

	for rows.Next() {
		val := make([]sql.NullString, len(cols))
		values := make([]interface{}, len(cols))
		for i := range cols {
			values[i] = &val[i]
		}

		if err := rows.Scan(values...); err != nil {
			return err
		}

		// var stringResult []string

		columNames, _ := rows.Columns()

		for i, v := range val {
			rowResult[columNames[i]] = v.String
		}
	}

	stValue := reflect.ValueOf(fields).Elem()
	sType := stValue.Type()
	for i := 0; i < sType.NumField(); i++ {
		field := sType.Field(i)
		fieldInterface := stValue.Field(i)
		if value, ok := rowResult[strings.Split(field.Tag.Get("json"), ",")[0]]; ok {
			if value != nil {
				switch fieldInterface.Interface().(type) {
				case string:
					stValue.Field(i).Set(reflect.ValueOf(value.(string)))
				case *string:
					valueString := value.(string)
					stValue.Field(i).Set(reflect.ValueOf(&valueString))
				case int:
					intString := strings.Split(value.(string), ".")[0]
					intValue, _ := strconv.Atoi(intString)
					stValue.Field(i).Set(reflect.ValueOf(intValue))
				case *int:
					intString := strings.Split(value.(string), ".")[0]
					intValue, _ := strconv.Atoi(intString)
					stValue.Field(i).Set(reflect.ValueOf(&intValue))
				case float64:
					float64Value, _ := strconv.ParseFloat(value.(string), 64)
					stValue.Field(i).Set(reflect.ValueOf(float64Value))
				case *float64:
					float64Value, _ := strconv.ParseFloat(value.(string), 64)
					stValue.Field(i).Set(reflect.ValueOf(&float64Value))
				case bool:
					boolValue, _ := strconv.ParseBool(value.(string))
					stValue.Field(i).Set(reflect.ValueOf(boolValue))
				case *bool:
					boolValue, _ := strconv.ParseBool(value.(string))
					stValue.Field(i).Set(reflect.ValueOf(&boolValue))
				case time.Time:
					// timeNew, _ := GetTime(value.(string), "2006-01-02 15:04:05")
					timeNew, _ := GetTime(value.(string), "2006-01-02T15:04:05+07:00")
					stValue.Field(i).Set(reflect.ValueOf(timeNew))
				case *time.Time:
					// timeNew, _ := GetTime(value.(string), "2006-01-02 15:04:05")
					if value.(string) != "" {
						timeNew, _ := GetTime(value.(string), "2006-01-02T15:04:05+07:00")
						stValue.Field(i).Set(reflect.ValueOf(&timeNew))
					}
				default:
					stValue.Field(i).Set(reflect.ValueOf(nil))
				}
			} else {
				stValue.Field(i).Set(reflect.ValueOf(value))
			}
		}
	}
	return nil
}

func MysqlQueryParams(tableName, as string, name string, fields any, params QueryParams) error {
	q := ""
	value := make([]interface{}, 0)
	if tableName != "" {
		if params.Select != nil {
			q += "SELECT "
			for _, v := range *params.Select {
				q += v + ", "
			}
			q = q[:len(q)-2]
		} else if as != "" {
			q += "SELECT " + as + ".*"
		} else {
			q += "SELECT " + tableName + ".*"
		}
		q += " FROM " + tableName + " " + as
		w := ""
		l := ""
		s := ""
		g := ""
		limit := ""
		offset := ""

		valueOn := make([]interface{}, 0)
		valueWhere := make([]interface{}, 0)
		if params.Where != nil {
			w = " WHERE "
			if params.Where.Field != nil && params.Where.Value != nil && params.Where.Op != nil {
				if *params.Where.Op == "IS" {
					value := *params.Where.Value
					w += *params.Where.Table + "." + *params.Where.Field + " " + *params.Where.Op + " " + value.(string)
				} else {
					w += *params.Where.Table + "." + *params.Where.Field + " " + *params.Where.Op + " ? "
					valueWhere = append(valueWhere, *params.Where.Value)
				}
			}

			if params.Where.Or != nil {
				for _, v := range *params.Where.Or {
					w += " OR " + OrQuery(v, &valueWhere)
				}
			}

			if params.Where.And != nil {
				for _, v := range *params.Where.And {
					w += " AND " + AndQuery(v, &valueWhere)
				}
			}

		}

		if params.LeftJoin != nil {
			for _, v := range *params.LeftJoin {
				q += " LEFT JOIN " + v.Table + " ON " + v.On
				if v.Value != nil {
					valueOn = append(valueOn, *v.Value...)
				}
			}
		}

		if params.Sort != nil {
			s = " ORDER BY "
			for _, v := range *params.Sort {
				s += v.Table + "." + v.Field + " " + v.Order + ", "
			}
			s = s[:len(s)-2]
		} else if as != "" {
			s = " ORDER BY " + as + ".id ASC"
		} else {
			s = " ORDER BY " + tableName + ".id ASC"
		}

		if params.GroupBy != nil {
			g = " GROUP BY "
			for _, v := range *params.GroupBy {
				g += v.Table + "." + v.Field + ", "
			}
			g = g[:len(g)-2]
		}

		if params.Limit != nil {
			limit = " LIMIT " + strconv.Itoa(*params.Limit)
		}

		if params.Offset != nil {
			offset = " OFFSET " + strconv.Itoa(*params.Offset)
		}

		q += l + w + g + s + limit + offset
		value = append(value, valueOn...)
		value = append(value, valueWhere...)
	}

	fmt.Println(q)
	rows, err := DB.Query(q, value...)
	if err != nil {
		return err
	}

	defer rows.Close()

	cols, _ := rows.Columns()

	result := make([]map[string]interface{}, 0)

	for rows.Next() {
		rowResult := make(map[string]interface{})
		val := make([]sql.NullString, len(cols))
		values := make([]interface{}, len(cols))
		for i := range cols {
			values[i] = &val[i]
		}

		if err := rows.Scan(values...); err != nil {
			return err
		}

		// var stringResult []string

		columNames, _ := rows.Columns()

		for i, v := range val {
			rowResult[columNames[i]] = v.String
		}
		result = append(result, rowResult)
	}

	stValue := reflect.ValueOf(fields).Elem()
	sType := stValue.Type()
	// returnValue := stValue
	returnType := sType
	returnValue := reflect.New(sType)
	if sType.Kind() == reflect.Slice {
		returnValue = reflect.MakeSlice(sType, len(result), len(result))
		sType = sType.Elem()
	}

	for idx, v := range result {
		if returnType.Kind() == reflect.Slice {
			stValue = returnValue.Index(idx)
		}
		for i := 0; i < sType.NumField(); i++ {
			field := sType.Field(i)
			fieldInterface := stValue.Field(i)
			if value, ok := v[strings.Split(field.Tag.Get("json"), ",")[0]]; ok {
				if value != nil {
					switch fieldInterface.Interface().(type) {
					case string:
						stValue.Field(i).Set(reflect.ValueOf(value.(string)))
					case *string:
						valueString := value.(string)
						stValue.Field(i).Set(reflect.ValueOf(&valueString))
					case int:
						intString := strings.Split(value.(string), ".")[0]
						intValue, _ := strconv.Atoi(intString)
						stValue.Field(i).Set(reflect.ValueOf(intValue))
					case *int:
						intString := strings.Split(value.(string), ".")[0]
						intValue, _ := strconv.Atoi(intString)
						stValue.Field(i).Set(reflect.ValueOf(&intValue))
					case float64:
						float64Value, _ := strconv.ParseFloat(value.(string), 64)
						stValue.Field(i).Set(reflect.ValueOf(float64Value))
					case *float64:
						float64Value, _ := strconv.ParseFloat(value.(string), 64)
						stValue.Field(i).Set(reflect.ValueOf(&float64Value))
					case bool:
						boolValue, _ := strconv.ParseBool(value.(string))
						stValue.Field(i).Set(reflect.ValueOf(boolValue))
					case *bool:
						boolValue, _ := strconv.ParseBool(value.(string))
						stValue.Field(i).Set(reflect.ValueOf(&boolValue))
					case time.Time:
						// timeNew, _ := GetTime(value.(string), "2006-01-02 15:04:05")
						timeNew, _ := GetTime(value.(string), "2006-01-02T15:04:05+07:00")
						stValue.Field(i).Set(reflect.ValueOf(timeNew))
					case *time.Time:
						// timeNew, _ := GetTimevalue.(string), "2006-01-02 15:04:05")
						if value.(string) != "" {
							timeNew, _ := GetTime(value.(string), "2006-01-02T15:04:05+07:00")
							stValue.Field(i).Set(reflect.ValueOf(&timeNew))
						}
					default:
						stValue.Field(i).Set(reflect.ValueOf(nil))
					}
				} else {
					stValue.Field(i).Set(reflect.ValueOf(value))
				}
			}
		}
		returnValue.Index(idx).Set(stValue)
	}
	// fields = returnValue
	stValue = reflect.ValueOf(fields).Elem()
	stValue.Set(returnValue)
	return nil
}

func InsertOrUpdateStruct(tableName string, data any) error {
	// Extract the struct's fields and values d
	val := reflect.ValueOf(data).Elem()
	typ := val.Type()
	numFields := typ.NumField()

	query := fmt.Sprintf("INSERT INTO %s (", tableName)
	value := make([]interface{}, 0)
	values := "VALUES ("
	update := ""
	for i := 0; i < numFields; i++ {
		field := typ.Field(i)
		if val.FieldByName(field.Name).Kind() == reflect.Ptr && !val.FieldByName(field.Name).IsNil() || val.FieldByName(field.Name).Kind() != reflect.Ptr {
			columnName := strings.Split(field.Tag.Get("json"), ",")[0]
			query += "`" + columnName + "`,"
			values += "?,"
			if len(update) == 0 {
				update += "ON DUPLICATE KEY UPDATE `" + columnName + "` = ? "
			} else {
				update += ", `" + columnName + "` = ?"
			}
			value = append(value, val.FieldByName(field.Name).Interface())
		}
	}

	// Trim the trailing commas and close the parentheses
	query = query[:len(query)-1] + ")"
	values = values[:len(values)-1] + ")"
	value = append(value, value...)
	// Combine the query and values to form the complete INSERT statement
	stmt := query + " " + values + " " + update

	// util.PrettyPrint(stmt, "S")
	// util.PrettyPrint(value, "V")

	// Prepare and execute the statement with the struct's values
	returnData, err := DB.Exec(stmt, value...)
	if err != nil {
		return err
	}

	id, _ := returnData.LastInsertId()
	idInt := int(id)
	reflect.ValueOf(data).Elem().FieldByName("ID").Set(reflect.ValueOf(&idInt))

	return err
}

func OrQuery(data OrParams, value *[]any) (result string) {
	if data.Field != nil && data.Value != nil && data.Op != nil {
		if data.Table != nil {
			if *data.Op == "IS" {
				value := *data.Value
				result += *data.Table + "." + *data.Field + " " + *data.Op + " " + value.(string)
			} else {
				result += *data.Table + "." + *data.Field + " " + *data.Op + " ? "
				*value = append(*value, *data.Value)
			}
		} else {
			if *data.Op == "IS" {
				value := *data.Value
				result += *data.Field + " " + *data.Op + " " + value.(string)
			} else {
				result += *data.Field + " " + *data.Op + " ? "
				*value = append(*value, *data.Value)
			}
		}
	}
	if data.And != nil {
		for _, v := range *data.And {
			result += " AND " + AndQuery(v, value)
		}
	}

	return
}

func AndQuery(data AndParam, value *[]any) (result string) {
	if data.Field != nil && data.Value != nil && data.Op != nil {
		if data.Table != nil {
			if *data.Op == "IS" {
				value := *data.Value
				result += *data.Table + "." + *data.Field + " " + *data.Op + " " + value.(string)
			} else {
				result += *data.Table + "." + *data.Field + " " + *data.Op + " ?"
				*value = append(*value, *data.Value)
			}
		} else {
			if *data.Op == "IS" {
				value := *data.Value
				result += *data.Field + " " + *data.Op + " " + value.(string)
			} else {
				result += *data.Field + " " + *data.Op + " ? "
				*value = append(*value, *data.Value)
			}
		}
	}
	if data.Or != nil {
		for _, v := range *data.Or {
			result += " OR " + OrQuery(v, value)
		}
	}
	return
}

func GetTime(params ...string) (time.Time, error) {
	layoutFormat := "2006-01-02T15:04:05.000Z"

	if len(params) < 1 {
		return time.Time{}, errors.New("need minimal one argument")
	}

	if len(params) > 1 {
		layoutFormat = params[1]
	}

	result, err := time.Parse(layoutFormat, params[0])

	return result, err
}

type QueryParams struct {
	Select   *[]string   `json:"select"`
	Where    *WhereParam `json:"where"`
	LeftJoin *[]LeftJoin `json:"left_join"`
	Sort     *[]Sort     `json:"sort"`
	GroupBy  *[]GroupBy  `json:"group_by"`
	Limit    *int        `json:"limit"`
	Offset   *int        `json:"offset"`
}

// type WhereParams struct {
// 	Table *string
// 	Field *string
// 	Op    *string
// 	Value *string
// }

type OrParams struct {
	Table *string
	Field *string
	Op    *string
	Value *any
	And   *[]AndParam
}

type AndParam struct {
	Table *string
	Field *string
	Op    *string
	Value *any
	Or    *[]OrParams
}

type WhereParam struct {
	Table *string
	Field *string
	Op    *string
	Value *any
	Or    *[]OrParams
	And   *[]AndParam
}
type LeftJoin struct {
	Table string
	On    string
	Value *[]any
}

type Sort struct {
	Table string
	Field string
	Order string
}

type GroupBy struct {
	Table string
	Field string
}

type SqlDB struct {
	*sql.DB
	QueryString string
}
