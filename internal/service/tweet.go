package service

import (
	"encoding/json"
	"fmt"
	"github.com/nestorov88/ElonTheGreat/pkg/tweet"
	"github.com/pkg/errors"
	"os"
)

type tweetService struct {
	repo tweet.TweetRepository
}

func NewTweetService(r tweet.TweetRepository) tweet.TweetService {
	return &tweetService{repo: r}
}

func (ts *tweetService) Store(tweet *tweet.Tweet) error {
	return ts.repo.Store(tweet)
}

func (ts *tweetService) ByTimeOfTheDay() ([]tweet.ByTimeOfTheDayResult, error) {
	return ts.repo.ByTimeOfTheDay()
}

func (ts *tweetService) StatsByYear() ([]tweet.StatsByYearResult, error) {
	return ts.repo.StatsByYear()
}

func (ts *tweetService) ImportJson(filePath string) (bool, error) {

	fmt.Println("FiletoImport: ", filePath)

	if len(filePath) == 0 {
		filePath = "elonmusk.json"
	}
	// Open our jsonFile
	jsonFile, err := os.Open(filePath)
	// if we os.Open returns an error then handle it
	if err != nil {
		return false, errors.Wrap(err, "tweetService.ImportJson.OpenFile")
	}

	fmt.Println("Successfully Opened ", filePath)
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	var tweets []tweet.Tweet
	json.NewDecoder(jsonFile).Decode(&tweets)

	for _, v := range tweets {
		ts.repo.StoreOrUpdate(&v)
	}

	fmt.Println("Importer executed successfully.")

	return true, nil
}
