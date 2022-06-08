package pkg

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

type Respsone struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Router
func Router() *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins: []string{"http://localhost:8080/*"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Post("/items", func(w http.ResponseWriter, r *http.Request) {
		var item Item
		if err := render.Bind(r, &item); err != nil {
			log.Println(err)
			renderResponse(w, map[string]string{"message": "failed to bind"}, http.StatusBadRequest)
			return
		}
		v, err := createItem(r.Context(), &item)
		if err != nil {
			log.Println(err)
			renderResponse(w, map[string]string{"message": "failed to create item"}, http.StatusBadRequest)
			return
		}

		renderResponse(w, Respsone{
			Message: "Ok",
			Data:    v,
		}, http.StatusCreated)
	})
	r.Get("/items", func(w http.ResponseWriter, r *http.Request) {
		v, err := listItems(r.Context())
		if err != nil {
			log.Println(err)
			renderResponse(w, map[string]string{"message": "failed to get items"}, http.StatusInternalServerError)
			return

		}
		renderResponse(w, Respsone{
			Message: "Ok",
			Data:    v,
		}, http.StatusOK)
	})
	r.Get("/items/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		v, err := getItem(r.Context(), id)
		if err != nil {
			log.Println(err)
			renderResponse(w, map[string]string{"message": "failed to get item"}, http.StatusInternalServerError)
			return
		}
		renderResponse(w, Respsone{
			Message: "Ok",
			Data:    v,
		}, http.StatusOK)
	})
	r.Delete("/items/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		err := deleteItem(r.Context(), id)
		if err != nil {
			log.Println(err)
			renderResponse(w, map[string]string{"message": "failed to delete item"}, http.StatusInternalServerError)
			return
		}
		renderResponse(w, Respsone{
			Message: "Ok",
			Data:    id,
		}, http.StatusOK)
	})

	FileServer(r)
	return r
}

// FileServer is serving static files.
func FileServer(router *chi.Mux) {
	root := "./static"
	fs := http.FileServer(http.Dir(root))

	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(root + r.RequestURI); os.IsNotExist(err) {
			http.StripPrefix(r.RequestURI, fs).ServeHTTP(w, r)
		} else {
			fs.ServeHTTP(w, r)
		}
	})
}
