package main

import (
	"net/url"

	"github.com/ChimeraCoder/anaconda"
)

type TwitterClient struct {
	api *anaconda.TwitterApi
}

func (client TwitterClient) Search(text string) anaconda.SearchResponse {
	searchResults, _ := client.api.GetSearch(text, nil)

	return searchResults
}

func (client TwitterClient) tweet(statusID string, message string) {
	v := url.Values{}
	v.Set("in_reply_to_status_id", statusID)
	client.api.PostTweet(message, v)
}

func (client TwitterClient) retweet(id int64) {
	client.api.Retweet(id, false)
}
