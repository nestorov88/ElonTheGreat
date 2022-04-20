package main

import (
	"fmt"
	"github.com/nestorov88/ElonTheGreat/internal/repository/mongodb"
	s "github.com/nestorov88/ElonTheGreat/internal/service"
	"os"
)

func main() {

	mongoRepo, err := mongodb.NewMongoRepository("mongodb://localhost:27017/?directConnection=true&serverSelectionTimeoutMS=2000", "elon_tweets", 10)

	if err != nil {
		fmt.Println("Terminating import: URL", err.Error())
		return
	}

	service := s.NewTweetService(mongoRepo)

	_, err = service.ImportJson(os.Getenv("TWEET_FILE"))
	//_, err = service.ImportJson("/home/sinemetu/go/src/elon_the_great/test/elonmusk.json")

	if err != nil {
		fmt.Printf("Importer occured error %s \n", err.Error())
		return
	}

	fmt.Println("Importer executed successfully.")
}
