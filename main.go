package main

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/ChimeraCoder/anaconda"
)

func main() {

	runtime.GOMAXPROCS(2)

	var wg sync.WaitGroup
	wg.Add(2)

	configuration, err := getConfig("config.json")
	if err != nil {
		fmt.Println(err)
	}
	anaconda.SetConsumerKey(configuration.ConsumerKey)
	anaconda.SetConsumerSecret(configuration.ConsumerSecret)
	api := anaconda.NewTwitterApi(configuration.AccessToken, configuration.AccessKey)
	twitterClient := TwitterClient{api: api}

	//parallelism! this runs both the search for Income Inequality and Women's Right at the same time
	//this optimizes time taken for each cycle

	go func() {

		defer wg.Done()

		results := twitterClient.Search("Income Inequality")
		tweets := results.Statuses
		for i := 0; i < len(tweets); i++ {
			tweet := tweets[i]
			retweeted := tweet.Retweeted
			if !retweeted {
				username := tweet.User.ScreenName
				fmt.Println(username + " : " + tweet.Text + "\n")
				message := "@" + username + " You should vote for Bernie Sanders! He wants to end income equality!"
				twitterClient.tweet(tweet.IdStr, message)
				twitterClient.retweet(tweet.Id)
			}
		}

	}()

	go func() {

		defer wg.Done()

		results := twitterClient.Search("Women's Rights")
		tweets := results.Statuses
		for i := 0; i < len(tweets); i++ {
			tweet := tweets[i]
			retweeted := tweet.Retweeted
			if !retweeted {

				username := tweet.User.ScreenName
				fmt.Println(username + " : " + tweet.Text + "\n")
				message := "@" + username + " You should vote for Bernie Sanders! He supports Women's Rights!"
				twitterClient.tweet(tweet.IdStr, message)
				twitterClient.retweet(tweet.Id)

			}
		}

	}()

	wg.Wait()

}
