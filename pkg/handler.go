package pkg

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/middleware"
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

	r.Use(middleware.Logger)
	
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

	r.Post("/items", create)

	r.Get("/items", get)

	r.Delete("/items/{id}", delete)

	r.Patch("/items/{id}", update)
	
	r.Put("/items/{id}", restore)

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

// create item handler
func create(w http.ResponseWriter, r *http.Request) {
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
}

// get items handler
func get(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	v, err := listItems(r.Context(), status)
	if err != nil {
		log.Println(err)
		renderResponse(w, map[string]string{"message": "failed to get items"}, http.StatusInternalServerError)
		return

	}
	renderResponse(w, Respsone{
		Message: "Ok",
		Data:    v,
	}, http.StatusOK)
}

// delete item handler
func delete(w http.ResponseWriter, r *http.Request) {
	var comment DeletedComment
	if err := render.Bind(r, &comment); err != nil {
		log.Println(err)
		renderResponse(w, map[string]string{"message": "failed to bind"}, http.StatusBadRequest)
		return
	}

	id := chi.URLParam(r, "id")
	err := deleteItem(r.Context(), id, &comment)
	if err != nil {
		log.Println(err)
		renderResponse(w, map[string]string{"message": "failed to delete item"}, http.StatusInternalServerError)
		return
	}
	renderResponse(w, Respsone{
		Message: "Ok",
		Data:    id,
	}, http.StatusOK)
}

// update item handler
func update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var item Item
	if err := render.Bind(r, &item); err != nil {
		log.Println(err)
		renderResponse(w, map[string]string{"message": "failed to bind"}, http.StatusBadRequest)
		return
	}
	err := updateItem(r.Context(), id, &item)
	if err != nil {
		log.Println(err)
		renderResponse(w, map[string]string{"message": "failed to update item"}, http.StatusInternalServerError)
		return
	}
	renderResponse(w, Respsone{
		Message: "Ok",
		Data:    id,
	}, http.StatusOK)
}

// restore item handler
func restore(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := unArchiveItem(r.Context(), id)
	if err != nil {
		log.Println(err)
		renderResponse(w, map[string]string{"message": "failed to restore item"}, http.StatusInternalServerError)
		return
	}
	renderResponse(w, Respsone{
		Message: "Ok",
		Data:   id,
	}, http.StatusOK)
}