package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

var users = make(map[int]*User)

var idCounter int

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/users", func(w http.ResponseWriter, _ *http.Request) {
		values := make([]*User, 0, len(users))
		for _, v := range users {
			values = append(values, v)
		}

		err := json.NewEncoder(w).Encode(values)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})

	r.Post("/users", func(w http.ResponseWriter, r *http.Request) {
		u := &User{}
		err := json.NewDecoder(r.Body).Decode(u)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		addUser(u)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(u)
	})

	r.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id := -1
		id, _ = strconv.Atoi(idStr)
		log.Println(id)
		user := users[id]
		if user == nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		err := json.NewEncoder(w).Encode(user)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})

	addUser(&User{Name: "Jeremy"})
	addUser(&User{Name: "George"})

	http.ListenAndServe(":3000", r)
}

func addUser(u *User) {
	idCounter += 1
	u.Id = idCounter

	users[u.Id] = u
}
