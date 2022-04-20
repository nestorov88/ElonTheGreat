package tweet

import "go.mongodb.org/mongo-driver/bson"

// Tweet
type Tweet struct {
	ID             int64    `json:"id"`
	ConversationID string   `json:"conversation_id"`
	CreatedAt      int64    `json:"created_at"`
	Date           string   `json:"date"`
	Time           string   `json:"time"`
	Timezone       string   `json:"timezone"`
	UserID         int      `json:"user_id"`
	Username       string   `json:"username"`
	Name           string   `json:"name"`
	Place          string   `json:"place"`
	Tweet          string   `json:"tweet"`
	Mentions       []string `json:"mentions"`
	Urls           []string `json:"urls"`
	Photos         []string `json:"photos"`
	RepliesCount   int      `json:"replies_count"`
	RetweetsCount  int      `json:"retweets_count"`
	LikesCount     int      `json:"likes_count"`
	Hashtags       []string `json:"hashtags"`
	Cashtags       []string `json:"cashtags"`
	Link           string   `json:"link"`
	Retweet        bool     `json:"retweet"`
	QuoteURL       string   `json:"quote_url"`
	Video          int      `json:"video"`
	Near           string   `json:"near"`
	Geo            string   `json:"geo"`
	Source         string   `json:"source"`
	UserRtID       string   `json:"user_rt_id"`
	UserRt         string   `json:"user_rt"`
	RetweetID      string   `json:"retweet_id"`
	ReplyTo        []User   `json:"reply_to"`
	RetweetDate    string   `json:"retweet_date"`
	Translate      string   `json:"translate"`
	TransSrc       string   `json:"trans_src"`
	TransDest      string   `json:"trans_dest"`
}

//User
type User struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}

// ByTimeOfTheDayResult is result stuct for
// TweetRepository.ByTimeOfTheDay query
type ByTimeOfTheDayResult struct {
	Hours int `json:"hours" bson:"_id"`
	Count int `json:"count" bson:"count"`
}

// StatsByYearResult is result stuct for
// TweetRepository.StatsByYear query
type StatsByYearResult bson.M

// TweetService is service interface that wraps
// Tweets domain business logic (import, store, fetch)
type TweetService interface {
	Store(tweet *Tweet) error
	ByTimeOfTheDay() ([]ByTimeOfTheDayResult, error)
	StatsByYear() ([]StatsByYearResult, error)
	ImportJson(filePath string) (bool, error)
}

// TweetRepository is interface that wraps
// Tweets db query logic
type TweetRepository interface {
	Store(tweet *Tweet) error
	StoreOrUpdate(tweet *Tweet) error
	ByTimeOfTheDay() ([]ByTimeOfTheDayResult, error)
	StatsByYear() ([]StatsByYearResult, error)
}
