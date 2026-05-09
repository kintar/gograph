package main

import (
	"fmt"
	"net/http"

	"github.com/ozgurcd/gograph/testdata/fixture/internal/auth"
	"github.com/ozgurcd/gograph/testdata/fixture/internal/db"
)

func main() {
	repo := &db.PostgresRepo{}
	svc := auth.NewService(repo)

	http.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		name, err := svc.Login(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		fmt.Fprintf(w, "Welcome %s", name)
	})

	http.ListenAndServe(":8080", nil)
}
