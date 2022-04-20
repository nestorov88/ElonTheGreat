package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	h "github.com/nestorov88/ElonTheGreat/internal/http/rest"
	"github.com/nestorov88/ElonTheGreat/internal/repository/mongodb"
	s "github.com/nestorov88/ElonTheGreat/internal/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	mongoRepo, err := mongodb.NewMongoRepository(os.Getenv("MONGO_URI"), "elon_tweets", 10)

	if err != nil {
		fmt.Printf("Terminated by error %s \n", err.Error())
		return
	}

	service := s.NewTweetService(mongoRepo)

	//If not empty the service will start the importing function first
	initSeed := os.Getenv("INIT_SEED")

	if len(initSeed) > 0 {

		//Opening path to json file to import
		_, err = service.ImportJson(os.Getenv("TWEET_FILE"))

		if err != nil {
			fmt.Printf("Importer occured error %s \n", err.Error())
		}
	}

	handler := h.New(service)

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/api/charts/time", handler.ByTimeOfTheDay)
	r.Get("/api/charts/stats", handler.StatsByYear)

	errs := make(chan error, 2)
	go func() {

		port := os.Getenv("PORT")

		if port == "" {
			port = "9000"
		}

		fmt.Println("Listening on port : ", port)
		errs <- http.ListenAndServe(":"+port, r)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("Terminated %s", <-errs)
}
