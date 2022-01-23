package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type Photo struct {
	Name string  `json:"photoname`
	Tags string `json:tags`
}

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "1122334455"
	DB_NAME     = "photoDB"
)

// подключение к БД
func setupDB() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Println(err)
	}

	return db
}
// Функция ищет по тегу в БД данные о картинке и возвращящет ее название со всеми тегами
func find(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	tag := r.URL.Query().Get("tag")

	rows, err := db.Query(fmt.Sprintf(`select name, array_agg(tag)  
	from phototags pt
	left join photos ph
	on pt.photoid = ph.id
	left join tags t
	on pt.tagid = t.tag_id
	group by name
	having '%s' = ANY (array_agg(tag))`, tag))
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	photos := []Photo{}

	for rows.Next() {
		p := Photo{}
		err := rows.Scan(&p.Name, &p.Tags)
		if err != nil {
			fmt.Println(err)
			continue
		}
		photos = append(photos, p)
	}

	json.NewEncoder(w).Encode(photos)

}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/find", find)

	log.Println("Server launch on localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)

}
