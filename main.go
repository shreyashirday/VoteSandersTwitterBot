package main

import (
	"fmt"
	"runtime"
	"strings"
	"sync"

	"github.com/ChimeraCoder/anaconda"
	"github.com/peterjliu/metamind"
)

func main() {

	runtime.GOMAXPROCS(2)

	var wg sync.WaitGroup
	wg.Add(2)

	configuration, err := getConfig("config.json")
	if err != nil {
		fmt.Println("config error")
		fmt.Println(err)
	} else {

		fmt.Println("good to go")

	}

	metamindClient := metamind.NewSentimentClient(configuration.MetamindKey)

	anaconda.SetConsumerKey(configuration.ConsumerKey)
	anaconda.SetConsumerSecret(configuration.ConsumerSecret)
	api := anaconda.NewTwitterApi(configuration.AccessToken, configuration.AccessKey)
	twitterClient := TwitterClient{api: api}

	sentiments := make(chan *metamind.SentimentResp)
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
				go metamind.GetSentiment(metamindClient, tweet.Text, sentiments)
			}
		}

		for i := 0; i < len(tweets); i++ {

			r := <-sentiments
			if r.Success {
				predictions := r.Prediction.Predictions
				var posValue float32
				var neuValue float32
				var negValue float32

				for i := 0; i < len(predictions); i++ {
					pred := predictions[i]
					if strings.EqualFold(pred.ClassName, "positive") {
						posValue = pred.Prob
					} else if strings.EqualFold(pred.ClassName, "negative") {
						negValue = pred.Prob
					} else {
						neuValue = pred.Prob
					}
				}

				if posValue > negValue && neuValue > negValue {
					tweet := tweets[i]
					username := tweet.User.ScreenName
					message := "@" + username + " You should vote for Bernie Sanders! He supports Women's Rights!"
					twitterClient.tweet(tweet.IdStr, message)
					twitterClient.retweet(tweet.Id)
				}
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

				go metamind.GetSentiment(metamindClient, tweet.Text, sentiments)
			}
		}

		for i := 0; i < len(tweets); i++ {

			r := <-sentiments
			if r.Success {
				predictions := r.Prediction.Predictions
				var posValue float32
				var neuValue float32
				var negValue float32

				for i := 0; i < len(predictions); i++ {
					pred := predictions[i]
					if strings.EqualFold(pred.ClassName, "positive") {
						posValue = pred.Prob
					} else if strings.EqualFold(pred.ClassName, "negative") {
						negValue = pred.Prob
					} else {
						neuValue = pred.Prob
					}
				}

				if posValue > negValue && neuValue > negValue {
					tweet := tweets[i]
					username := tweet.User.ScreenName
					message := "@" + username + " You should vote for Bernie Sanders! He supports Women's Rights!"
					twitterClient.tweet(tweet.IdStr, message)
					twitterClient.retweet(tweet.Id)
				}

			}

		}

	}()

	wg.Wait()

}
