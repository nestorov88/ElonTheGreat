package rest

import (
	"encoding/json"
	"github.com/nestorov88/ElonTheGreat/pkg/tweet"
	"github.com/pkg/errors"
	"net/http"
)

//Implements http handler func
type handler struct {
	s tweet.TweetService
}

//New export creating of handler
func New(s tweet.TweetService) *handler {
	return &handler{s}
}

// ByTimeOfTheDay query and return json with all tweets grouped by
// time of the day they were tweeted
func (h *handler) ByTimeOfTheDay(w http.ResponseWriter, r *http.Request) {

	tweetsByTime, err := h.s.ByTimeOfTheDay()

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	rawMsg, err := json.Marshal(tweetsByTime)

	if err != nil {
		errors.Wrap(err, "handler.ByTimeOfTheDay.Encode")
	}

	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

	w.Write(rawMsg)
}

// StatsByYear query and return json with count of tweets grouped by
// year and their likes and retweets
func (h *handler) StatsByYear(w http.ResponseWriter, r *http.Request) {

	tweetsStatsByYear, err := h.s.StatsByYear()

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	rawMsg, err := json.Marshal(tweetsStatsByYear)

	if err != nil {
		errors.Wrap(err, "handler.StatsByYear.Encode")
	}

	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

	w.Write(rawMsg)
}
