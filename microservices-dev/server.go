package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

//
// `memcached` 的 API 非常简单但又强大，有两种方法供用户使用 Get 与 Set
// 有一点值得注意的是使用这俩方法前，需要转换数据为 `[]byte`
type Name struct {
	NConst    string `json:"nconst"`
	Name      string `json:"name"`
	BirthYear string `json:"birthYear"`
	DeathYear string `json:"deathYear"`
}

type Error struct {
	Message string `json:"error"`
}

// 工作流总是一致的：
// 1. 从 memcached 中获取值，如果存在则返回
// 2. 如果不存在，查询原始数据，并存储于 memcached 中
func main() {

	// 数据库
	db, err := NewPostgreSQL()
	if err != nil {
		log.Fatalf("Could not initialize Database connection %s", err)
	}
	defer db.Close()

	// Memcached
	mc, err := NewMemcached()
	if err != nil {
		log.Fatalf("Could not initialize Memcached client %s", err)
	}

	router := mux.NewRouter()

	renderJSON := func(w http.ResponseWriter, val interface{}, statusCode int) {
		w.WriteHeader(statusCode)
		_ = json.NewEncoder(w).Encode(val)
	}

	router.HandleFunc("/names/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		val, err := mc.GetName(id)
		if err != nil {
			renderJSON(w, &val, http.StatusOK)
			return
		}

		name, err := db.FindByNConst(id)
		if err != nil {
			renderJSON(w, &Error{Message: err.Error()}, http.StatusInternalServerError)
			return
		}

		_ = mc.SetName(name)

		renderJSON(w, &name, http.StatusOK)
	})

	fmt.Println("Starting server :8080")

	srv := &http.Server{
		Handler:      router,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
