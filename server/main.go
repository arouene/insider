package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"insider/game"
	"insider/middlewares"
	"insider/session"

	. "insider/utils"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

const (
	configFileName = "insider"
	configFileType = "toml"
)

var (
	debug   bool
	timeout int
)

func init() {
	// Configuration initialisation
	viper.SetDefault("debug", false)
	viper.SetDefault("timeout", 240) // 4 minutes

	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFileType)

	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Config file %s cannot be read\n", configFileName+"."+configFileType)
	} else {
		log.Printf("Config file used: %s\n", viper.ConfigFileUsed())
	}

	viper.SetEnvPrefix("insider")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", " ", "_"))
	viper.AutomaticEnv()

	debug = viper.GetBool("debug")
	timeout = viper.GetInt("timeout")
	words := viper.GetStringSlice("words")

	game.InitWords(words)

	// Debug mode
	if !debug {
	}

	// Initialise randomness
	rand.Seed(time.Now().Unix())
}

func main() {
	r := mux.NewRouter()

	// CORS handling
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "HEAD"})
	origins := handlers.AllowedOrigins([]string{"*"})
	r.Use(handlers.CORS(headers, methods, origins))

	r.Use(middlewares.Log)
	//r.SkipClean(true)

	r.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"err": "Page not found"}`, http.StatusNotFound)
	}))
	r.Handle("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"msg": "pong"}`)
	}))

	game := r.PathPrefix("/game").Subrouter()
	game.Use(session.Inject)
	
	// Create a game or join a game
	game.HandleFunc("/{name:[a-zA-Z0-9]+}/join", join).
		Queries("playerName", "{playerName}").
		Queries("uuid", "{uuid:[a-z0-9-]*}").
		Methods("GET")
	// Transition GameState from CREATED to SETUP
	game.HandleFunc("/{name:[a-zA-Z0-9]+}/ready", ready).
		Queries("uuid", "{uuid:[a-z0-9-]*}").
		Methods("GET")
	// Transition GameState from SETUP to STARTED
	game.HandleFunc("/{name:[a-zA-Z0-9]+}/start", start).
		Queries("uuid", "{uuid:[a-z0-9-]*}").
		Methods("GET")
	// Transition GameState from STARTED to STOPPED
	game.HandleFunc("/{name:[a-zA-Z0-9]+}/stop", stop).
		Queries("uuid", "{uuid:[a-z0-9-]*}").
		Methods("GET")
	// Transition GameState from STOPPED to RESOLVED
	game.HandleFunc("/{name:[a-zA-Z0-9]+}/resolve", resolve).
		Queries("uuid", "{uuid:[a-z0-9-]*}").
		Methods("GET")
	// Optionally reset to CREATED
	game.HandleFunc("/{name:[a-zA-Z0-9]+}/reset", reset).
		Queries("uuid", "{uuid:[a-z0-9-]*}").
		Methods("GET")
	// Get game state
	game.HandleFunc("/{name:[a-zA-Z0-9]+}/state", state).
		Queries("uuid", "{uuid:[a-z0-9-]*}").
		Methods("GET")
	// Choose a new word when GameState in SETUP
	game.HandleFunc("/{name:[a-zA-Z0-9]+}/word", newWord).
		Queries("uuid", "{uuid:[a-z0-9-]*}").
		Methods("GET")

	// Launch server
	log.Fatal(http.ListenAndServe(":3000", r))
}

func join(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value(session.Key).(session.SessionContext)

	vars := mux.Vars(r)
	playerName := vars["playerName"]
	id := ctx.ID

	if ctx.Session.Game.Phase != game.CREATED {
		if ctx.IsIdentified() {
			JSON(w, Struct{}, http.StatusOK)
		} else {
			JSON(w, Struct{"err": "Game not in CREATED phase"}, http.StatusBadRequest)
		}
		return
	}

	// Creating player
	newPlayer := game.NewPlayer(playerName, id)
	ctx.Session.Players.AddPlayer(newPlayer)
	JSON(w, Struct{"name": newPlayer.Name, "id": newPlayer.Id}, http.StatusCreated)
	return
}

func ready(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value(session.Key).(session.SessionContext)

	if !ctx.IsIdentified() {
		JSON(w, Struct{}, http.StatusUnauthorized)
		return
	}

	if ctx.Session.Game.Phase != game.CREATED {
		JSON(w, Struct{"err": "Game not in CREATED phase"}, http.StatusBadRequest)
		return
	}

	if err := ctx.Session.Players.SetRandomRoles(); err != nil {
		JSON(w, Struct{"err": err.Error()}, http.StatusBadRequest)
		return
	}

	if err := ctx.Session.Game.SetPhase(game.SETUP); err != nil {
		JSON(w, Struct{"err": err.Error()}, http.StatusBadRequest)
		return
	}
}

func start(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value(session.Key).(session.SessionContext)

	if !ctx.IsMaster() {
		JSON(w, Struct{}, http.StatusUnauthorized)
		return
	}

	if ctx.Session.Game.Phase != game.SETUP {
		JSON(w, Struct{"err": "Game not in SETUP phase"}, http.StatusBadRequest)
		return
	}

	ctx.Session.Game.SetStartTime()

	if err := ctx.Session.Game.SetPhase(game.STARTED); err != nil {
		JSON(w, Struct{"err": err.Error()}, http.StatusBadRequest)
		return
	}
}

func stop(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value(session.Key).(session.SessionContext)

	if !ctx.IsMaster() {
		JSON(w, Struct{}, http.StatusUnauthorized)
		return
	}

	if err := ctx.Session.Game.SetPhase(game.STOPPED); err != nil {
		JSON(w, Struct{"err": err.Error()}, http.StatusBadRequest)
		return
	}
}

func resolve(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value(session.Key).(session.SessionContext)

	if !ctx.IsMaster() {
		JSON(w, Struct{}, http.StatusUnauthorized)
		return
	}

	if err := ctx.Session.Game.SetPhase(game.RESOLVED); err != nil {
		JSON(w, Struct{"err": err.Error()}, http.StatusBadRequest)
		return
	}
}

func reset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value(session.Key).(session.SessionContext)

	if !ctx.IsIdentified() {
		JSON(w, Struct{}, http.StatusUnauthorized)
		return
	}

	if err := ctx.Session.Game.SetPhase(game.CREATED); err != nil {
		JSON(w, Struct{"err": err.Error()}, http.StatusBadRequest)
		return
	}

	ctx.Session.Game.SetNewWord()
}

func state(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value(session.Key).(session.SessionContext)

	p := ctx.BuildPresenter()
	if ctx.IsMaster() || ctx.IsInsider() {
		p.WithWord()
	}

	if ctx.Session.Game.Phase == game.RESOLVED {
		p.WithPlayersRoles().WithWord()
	}

	JSON(w, Struct{"state": p.Build()}, 200)
}

func newWord(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value(session.Key).(session.SessionContext)

	if !ctx.IsMaster() {
		JSON(w, Struct{}, http.StatusUnauthorized)
		return
	}

	if ctx.Session.Game.Phase != game.SETUP {
		JSON(w, Struct{"err": "Game not in SETUP phase"}, http.StatusBadRequest)
		return
	}

	newWord := ctx.Session.Game.SetNewWord()

	JSON(w, Struct{"word": newWord}, http.StatusOK)
}
