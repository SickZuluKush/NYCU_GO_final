package db

import (
	"fmt"
	"log"
)


func Reset_record(filename string) {
	q := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", 
		schema.comment.table_name, schema.comment.filename) 
	if _, err := db.Exec(q, filename); err != nil {
		log.Print(err)
	}
}

func Add_record(filename string, comment string) {
	q := fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES (?, ?)", 
		schema.comment.table_name, schema.comment.filename, schema.comment.comment)
	
	if _, err := db.Exec(q, filename, comment); err != nil {
		log.Print(err)
	}
}

func Get_record(filename string) []string {
	q := fmt.Sprintf("SELECT %s FROM %s WHERE %s = ?", 
		schema.comment.comment, schema.comment.table_name, schema.comment.filename)
	
	rows, err := db.Query(q, filename)
	if err != nil {
		log.Fatalf("query %s failed: %s", filename, err)
		return []string{}
	}
	
	defer rows.Close()
	
	var results []string
	
	for rows.Next() {
		var r string
		if err := rows.Scan(&r); err != nil {
			log.Fatalf("query %s failed: %s", filename, err)
			return []string{}
		}
		
		results = append(results, r)
	}
	
	if err := rows.Err(); err != nil {
		log.Fatalf("query %s failed: %s", filename, err)
		return []string{}
	}
	
	return results
}

/*
func Debug() {
	q := fmt.Sprintf("SELECT * FROM %s", schema.comment.table_name)
	rows, err := db.Query(q)
	if err != nil {
		log.Fatal(err)
	}
	
	defer rows.Close()
	
	for rows.Next() {
		var (
			filename string
			comment string
		)
		
		if err := rows.Scan(&filename, &comment); err != nil {
			log.Fatal(err)
		}
		log.Printf("file: %s, comment: %s", filename, comment)
	}
	log.Print()
	
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
*/
