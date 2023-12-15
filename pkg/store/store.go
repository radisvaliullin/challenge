package store

import (
	"log"
	"net/http"
	"time"

	"github.com/radisvaliullin/challenge/pkg/api"
	"github.com/radisvaliullin/challenge/pkg/storage"
)

type Config struct {
	Addr string
}

type Store struct {
	conf Config
}

func New(conf Config) *Store {
	s := &Store{
		conf: conf,
	}
	return s
}

func (s *Store) Run() error {

	// run dependencies
	strgConf := storage.Config{
		Fixtures: storage.GetTestConfigFixtures(),
	}
	storage, err := storage.New(strgConf)
	if err != nil {
		log.Printf("store: run: new storage, err: %v", err)
		return err
	}

	// set api
	apiMux := api.New(storage).Muxer()

	// set http server
	srv := &http.Server{
		Addr:         s.conf.Addr,
		Handler:      apiMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return srv.ListenAndServe()
}
