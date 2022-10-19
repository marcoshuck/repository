package main

import (
	"context"
	"github.com/marcoshuck/repository"
	"gorm.io/gorm"
)

type Test struct {
	gorm.Model
	FirstName string
	LastName  string
}

func main() {
	r := repository.NewRepositorySQL[Test]()

	r.Create(context.Background(), Test{
		FirstName: "Test",
		LastName:  "Test",
	})

	r.Get(context.Background(), 0)
}
