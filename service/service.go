package service

import "os"

type Base interface {
	Run()
}

func Create() Base {
	if os.Getenv("MODE") == "http" {
		return &Http{}
	}
	return &Cli{}
}
