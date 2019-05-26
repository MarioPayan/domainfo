package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func startApi(port string) {
	router := Routes()
	showRoutes(router)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func Routes() *chi.Mux {
	router := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "user-auth"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.DefaultCompress,
		middleware.RedirectSlashes,
		middleware.Recoverer,
		cors.Handler,
	)
	router.Route("/v1", func(r chi.Router) {
		r.Mount("/api/domain", domainRoutes())
	})
	return router
}

func showRoutes(router *chi.Mux) {
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s\n", method, route)
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		log.Panicf("Logging err: %s\n", err.Error())
	}
}

func domainRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/{domain}", GetDomain)
	router.Get("/list", GetDomains)
	router.Delete("/{domain}", DeleteDomain)
	return router
}

func GetDomain(w http.ResponseWriter, r *http.Request) {
	domainID := chi.URLParam(r, "domain")
	domain, exist := getDomain(domainID)
	oneHourAfterUpdate := domain.UpdatedAt.Add(time.Hour * 1)
	if !exist {
		domainData := getData(domainID)
		if domainData.Status != READY {
			if domainData.Status == IN_PROGRESS {
				domain = createDomain(domainData)
				render.Status(r, 206)
			}
			if domainData.Status == ERROR {
				render.Status(r, 210)
			}
			if domainData.Status == DNS {
				render.Status(r, 211)
			}
		} else {
			for _, server := range domainData.Servers {
				if server.Status != Ready {
					render.Status(r, 206)
				}
			}
			domain = createDomain(domainData)
		}
	} else if domain.Status == IN_PROGRESS || oneHourAfterUpdate.Before(time.Now().UTC()) {
		render.Status(r, 206)
		domainData := getData(domainID)
		domain = updateDomain(domainData)
	}
	render.JSON(w, r, domain)
}

func GetDomains(w http.ResponseWriter, r *http.Request) {
	domains := getDomains()
	render.JSON(w, r, domains)
}

func DeleteDomain(w http.ResponseWriter, r *http.Request) {
	domainID := chi.URLParam(r, "domain")
	deleteDomain(domainID)
}
