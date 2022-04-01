package main

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/mwolf1989/sqlexec/util"
	"os"
	"runtime"
	"time"
)

func main() {

	argsWithoutProg := os.Args[1:]
	kind := argsWithoutProg[0]
	query := argsWithoutProg[1]
	configname := argsWithoutProg[2]
	if configname == "" {
		configname = "app"
	}
	config, err := util.LoadConfig(".", configname)
	if err != nil {
		fmt.Printf("Error loading config: %s\n", err)
		os.Exit(2)
	}
	db, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		fmt.Printf("Error connecting database: %s\n", err)
		os.Exit(14)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Printf("Error closing database connection: %s\n", err)
			os.Exit(14)
		}
	}(db)
	fmt.Printf("Executing query: %s\n", query)
	if kind == "exec" {
		res, err := db.Exec(query)
		if err != nil {
			fmt.Printf("Error executing query: %s\n", err)
			os.Exit(1)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			fmt.Printf("Error getting rows affected: %s\n", err)
			os.Exit(1)
		}
		lastId, err := res.LastInsertId()
		if err != nil {
			fmt.Printf("Error getting last insert id: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("Rows Affected: %d, last inserted id: %d", affected, lastId)
		os.Exit(0)
	} else if kind == "query" {
		rows, err := db.Query(query)
		if err != nil {
			fmt.Printf("Error executing query: %s\n", err)
			os.Exit(1)
		}
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {
				fmt.Printf("Error closing rows: %s\n", err)
				os.Exit(1)
			}
		}(rows)
		printTable(rows)
		os.Exit(0)
	} else {
		fmt.Printf("Unknown type: %s\n", kind)
		os.Exit(1)
	}

}

func printTable(rows *sql.Rows) {

	pr := func(t interface{}) (r string) {
		r = "\\N"
		switch v := t.(type) {
		case *sql.NullBool:
			if v.Valid {
				r = fmt.Sprintf("%v", v.Bool)
			}
		case *sql.NullString:
			if v.Valid {
				r = v.String
			}
		case *sql.NullInt64:
			if v.Valid {
				r = fmt.Sprintf("%6d", v.Int64)
			}
		case *sql.NullFloat64:
			if v.Valid {
				r = fmt.Sprintf("%.2f", v.Float64)
			}
		case *time.Time:
			if v.Year() > 1900 {
				r = v.Format("_2 Jan 2006")
			}
		default:
			r = fmt.Sprintf("%#v", t)
		}
		return
	}

	c, _ := rows.Columns()
	n := len(c)

	// print labels
	for i := 0; i < n; i++ {
		if len(c[i]) > 1 && c[i][1] == ':' {
			fmt.Print(c[i][2:], "\t")
		} else {
			fmt.Print(c[i], "\t")
		}
	}
	fmt.Print("\n\n")

	// print data
	var field []interface{}
	for i := 0; i < n; i++ {
		switch {
		case c[i][:2] == "b:":
			field = append(field, new(sql.NullBool))
		case c[i][:2] == "f:":
			field = append(field, new(sql.NullFloat64))
		case c[i][:2] == "i:":
			field = append(field, new(sql.NullInt64))
		case c[i][:2] == "s:":
			field = append(field, new(sql.NullString))
		case c[i][:2] == "t:":
			field = append(field, new(time.Time))
		default:
			field = append(field, new(sql.NullString))
		}
	}
	for rows.Next() {
		checkErr(rows.Scan(field...))
		for i := 0; i < n; i++ {
			fmt.Print(pr(field[i]), "\t")
		}
		fmt.Println()
	}
	fmt.Println()
}

func checkErr(err error) {
	if err != nil {
		_, filename, lineno, ok := runtime.Caller(1)
		if ok {
			fmt.Fprintf(os.Stderr, "%v:%v: %v\n", filename,
				lineno, err)
		}
		panic(err)
	}
}
