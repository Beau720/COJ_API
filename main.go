package main

import (
	"log"

	"COJ_API/http"
	"COJ_API/service/database"
	"COJ_API/service/form"
	"COJ_API/service/user"

	"github.com/kardianos/service"
)

var logger service.Logger

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	//Connect to the DB
	db_config := &database.Config{
		Username:       "root",
		Password:       "FireStore@321$",
		ConnectionType: "tcp",
		Host:           "127.0.0.1",
		Port:           3306,
		Name:           "coj",
	}

	http_config := &http.Config{
		Host:        "",
		Port:        "8080",
		UserService: user.NewUser(db_config),
		FormService: form.NewForm(db_config),
	}

	http.Start(http_config)
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func main() {

	svcConfig := &service.Config{
		Name:        "COJ_API",
		DisplayName: "COJ API",
		Description: "This is just an API for the COJ demo that listens to port 8080.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)

	}

	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
