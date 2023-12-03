package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/Foody-App-Tech/Main-server/config"
	"github.com/Foody-App-Tech/Main-server/internal/businesses"
	db "github.com/Foody-App-Tech/Main-server/internal/db/sqlc"
	mw "github.com/Foody-App-Tech/Main-server/internal/global_middlewares"
	"github.com/Foody-App-Tech/Main-server/internal/login"
	"github.com/Foody-App-Tech/Main-server/internal/users"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {

	// setup exit code for shutdown
	var exitCode int
	defer func() {
		os.Exit(exitCode)
	}()

	// load config
	env, err := config.LoadEnv("./")
	if err != nil {
		log.Printf("error --> %v", err)
		exitCode = 1
		return
	}

	// run the server
	err = run(env)
	if err != nil {
		log.Printf("error --> %v", err)
		exitCode = 1
		return
	}

}

func run(env config.EnvVariables) (err error) {

	srv, err := buildServer(env)
	if err != nil {
		return
	}

	// start server
	err = srv.ListenAndServe()
	if err != nil {
		log.Println("The authentication server could not start !")
		return
	}

	return
}

func buildServer(env config.EnvVariables) (srv *http.Server, err error) {
	// ----- connect to db
	conn, err := sql.Open(env.DB_DRIVER, env.DB_SOURCE)
	log.Println("Connected to foody_db successfully")
	if err != nil {
		log.Println("DB could not start")
		return
	}

	// TODO: run db migrations

	store := db.NewStore(conn)

	// ----- create *http.Server
	log.Println("Starting server on port", env.PORT)

	r := chi.NewRouter()

	// server rules
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-type", "X-CSRF-Token", "foody-api-key"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
	}))

	// midlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	globalHelpers := mw.NewGlobalMiddlewareService(env.ACCESS_TOKEN_KEY, env.ACCESS_TOKEN_TIME, env.REFRESH_TOKEN_KEY, env.REFRESH_TOKEN_TIME, env.API_KEY, env.SECURE_COOKIES, store)

	// auth domain
	loginService := login.NewLoginService(store, env)
	googleOauthService := login.NewGoogleOauthService(store, env)
	facebookOauthService := login.NewFacebookOauthService(store, env)
	loginController := login.NewLoginController(googleOauthService, facebookOauthService, loginService, globalHelpers)
	login.AddLoginRoutes(r, loginController)

	// users domain
	userService := users.NewUserService(store, env)
	userController := users.NewUserController(userService, globalHelpers)
	users.AddUserRouter(r, userController)

	// businesses domain
	businessesService := businesses.NewBusinessesService(store, env)
	businessesController := businesses.NewBusinessesController(businessesService, globalHelpers)
	businesses.AddBusinessRouter(r, businessesController)

	// health route
	r.Get("/health-check", func(w http.ResponseWriter, r *http.Request) {
		response := config.ClientResponse{
			Rsp: "server up !!!",
		}
		config.WriteResponse(w, 200, response)
	})

	srv = &http.Server{
		Addr:    fmt.Sprintf(":%s", env.PORT),
		Handler: r,
	}

	return
}
