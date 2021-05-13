package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Client Client `yaml:"client"`
	Rules  []Rule `yaml:"rules"`
}

type Client struct {
	ConsumerKey    string `yaml:"consumer_key"`
	ConsumerSecret string `yaml:"consumer_secret"`
	AccessToken    string `yaml:"access_token"`
	AccessSecret   string `yaml:"access_secret"`
}

type Options struct {
	ExcludeReplies  bool `yaml:"excludeReplies"`
	IncludeRetweets bool `yaml:"includeRetweets"`
}

type Actions map[string]interface{}

type Rule struct {
	Account  string  `yaml:"account"`
	Likes    string  `yaml:"likes"`
	Retweets string  `yaml:"retweets"`
	Options  Options `yaml:"options"`
	Actions  Actions `yaml:"actions"`
}

func main() {
	var (
		configPath string
		config     Config
		wg         sync.WaitGroup
	)

	flag.StringVar(&configPath, "config", "", "")
	flag.Parse()

	file, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		panic(err)
	}

	c := oauth1.NewConfig(config.Client.ConsumerKey, config.Client.ConsumerSecret)
	t := oauth1.NewToken(config.Client.AccessToken, config.Client.AccessSecret)

	httpClient := c.Client(oauth1.NoContext, t)
	client := twitter.NewClient(httpClient)

	wg.Add(len(config.Rules))

	for _, rule := range config.Rules {
		go func(rule Rule) {
			defer wg.Done()
			tweets, _, err := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
				ScreenName:      rule.Account,
				ExcludeReplies:  &rule.Options.ExcludeReplies,
				IncludeRetweets: &rule.Options.IncludeRetweets,
				Count:           100,
			})
			if err != nil {
				log.Println(err)
				return
			}

			for _, tweet := range tweets {
				log.Printf("account=%s id=%s reetweets=%d favorites=%d", rule.Account, tweet.IDStr, tweet.RetweetCount, tweet.FavoriteCount)

				retweetOperator := string(rule.Retweets[0])
				retweetCount, err := strconv.Atoi(string(rule.Retweets[1:]))
				if err != nil {
					log.Println(err)
					break
				}
				likeOperator := string(rule.Likes[0])
				likeCount, err := strconv.Atoi(string(rule.Likes[1:]))
				if err != nil {
					log.Println(err)
					break
				}

				if retweetOperator == ">" && likeOperator == ">" {
					if tweet.RetweetCount >= retweetCount && tweet.FavoriteCount >= likeCount {
						if _, ok := config.Rules[0].Actions["favorite"]; ok {
							_, _, err := client.Favorites.Create(&twitter.FavoriteCreateParams{
								ID: tweet.ID,
							})
							if err != nil {
								log.Println(err)
							} else {
								log.Printf("fav=%s", tweet.IDStr)
							}
						}

						if _, ok := config.Rules[0].Actions["retweet"]; ok {
							_, _, err := client.Statuses.Retweet(tweet.ID, &twitter.StatusRetweetParams{})
							if err != nil {
								log.Println(err)
							} else {
								log.Printf("retweet=%s", tweet.IDStr)
							}
						}
					}
				}
			}
		}(rule)
	}

	wg.Wait()

}
