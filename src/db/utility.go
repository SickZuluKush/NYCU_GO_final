package db

import (
	"os"
	"log"
	"fmt"
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db_folder string
	db_name string
	db *sql.DB
	schema db_tables
)

type db_tables struct {
	comment comment_schema
}

type comment_schema struct {
	table_name string
	filename string
	comment string
}


func init() {
	var ok bool
	
	if err := godotenv.Load(); err != nil {
		log.Print("load .env error")
	}
	
	if db_folder, ok = os.LookupEnv("DATABASE_FOLDER_NAME"); !ok {
		log.Print(".env error: db folder")
		db_folder = "db_folder"
	}
	if db_name, ok = os.LookupEnv("DATABASE_NAME"); !ok {
		log.Print(".env error: db name")
		db_name = "db_name"
	}
	
	comment_table_name, ok := os.LookupEnv("COMMENT_TABLE_NAME")
	if !ok {
		log.Print(".env error: comment table name")
		comment_table_name = "comment_table"
	}
	
	file_col_name, ok := os.LookupEnv("FILE_COL_NAME")
	if !ok {
		log.Print(".env error: file column name")
		file_col_name = "file_col"
	}
	
	comment_col_name, ok := os.LookupEnv("COMMENT_COL_NAME")
	if !ok {
		log.Print(".env error: comment column name")
		comment_col_name = "comment_col"
	}
	
	if err := os.MkdirAll(db_folder, os.ModePerm); err != nil {
		log.Print("mkdir err: db folder")
	}
	
	schema = db_tables{
		comment: comment_schema{
			comment_table_name, 
			file_col_name, 
			comment_col_name,
		},
	}
	
	init_db()
}

func init_db() {
	var err error
	db_path := db_folder + "/" + db_name
	
	db, err = sql.Open("sqlite3", db_path)
	if err != nil {
		db.Close()
		log.Fatal(err)
	}
	
	q := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s TEXT, %s TEXT)", 
		schema.comment.table_name, schema.comment.filename, schema.comment.comment)
	if _, err = db.Exec(q); err != nil {
		db.Close()
		log.Fatal(err)
	}
}
