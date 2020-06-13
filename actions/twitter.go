package actions

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"lickerbot/models"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop/v5"
	"github.com/sirupsen/logrus"
)

const (
	lickerbotUserID = 1269307350520868866
)

var twitterConsumerSecret = envy.Get("TWITTER_CONSUMER_SECRET", "")

var lickReactions = []string{
	"Liiiiick.",
	"That's one hell of a lick!",
	"Must taste good.",
}

// https://developer.twitter.com/en/docs/accounts-and-users/subscribe-account-activity/guides/account-activity-data-objects#tweet_create_events
type accountActivityEvent struct {
	ForUserID         int64   `json:"for_user_id,string"`
	TweetCreateEvents []tweet `json:"tweet_create_events"`
}

// https://developer.twitter.com/en/docs/tweets/data-dictionary/overview/Tweet-object
type tweet struct {
	ID                  int64  `json:"id"`
	InReplyToStatusID   int64  `json:"in_reply_to_status_id"`
	InReplyToUserID     int64  `json:"in_reply_to_user_id"`
	InReplyToScreenName string `json:"in_reply_to_screen_name"`
	Text                string `json:"text"`
	// https://developer.twitter.com/en/docs/tweets/data-dictionary/overview/entities-object
	Entities struct {
		UserMentions []userMention `json:"user_mentions"`
	} `json:"entities"`
	User struct {
		ScreenName string `json:"screen_name"`
	} `json:"user"`
}

// https://developer.twitter.com/en/docs/tweets/data-dictionary/overview/entities-object#mentions
type userMention struct {
	// ID of the mentioned user
	ID int64 `json:"id"`
}

// TwitterCRC implements https://developer.twitter.com/en/docs/accounts-and-users/subscribe-account-activity/guides/securing-webhooks
func TwitterSecurityCheck(c buffalo.Context) error {
	crcToken := c.Params().Get("crc_token")
	responseToken := hmac256(crcToken, twitterConsumerSecret)
	return c.Render(http.StatusOK, r.JSON(map[string]string{
		"response_token": "sha256=" + responseToken,
	}))
}

// TwitterWebhook default implementation.
func TwitterWebhook(c buffalo.Context) error {
	log := logrus.WithField("handler", "TwitterWebhook")

	event := &accountActivityEvent{}
	err := c.Bind(event)
	// if we fail to bind, assume this is a request for an account activity event
	// we don't care about. return a 200 to the webhook api so it doesn't get retried.
	// (though i believe on the free tier retries are disabled anyway.)
	if err != nil {
		log.WithError(err).Warn("error binding request")
		return c.Render(http.StatusOK, r.JSON(nil))
	}

	// we should only receive account activity for @lickerbot, return a 200 so it doesn't
	// get retried.
	if event.ForUserID != lickerbotUserID {
		log.Warnf("received webhook event for different account than lickerbot, %d", event.ForUserID)
		return c.Render(http.StatusOK, r.JSON(nil))
	}

	// pluck our actionable tweets from the event
	var actionableTweets []tweet
	for _, tweet := range event.TweetCreateEvents {
		if err := validate(&tweet); err != nil {
			log.Infof("received non-actionable tweet %d in webhook (%v)\n", tweet.ID, err)
		} else {
			actionableTweets = append(actionableTweets, tweet)
		}
	}

	// no actionable tweets, return a 200
	if len(actionableTweets) == 0 {
		log.Infof("no actionable tweets in webhook event")
		return c.Render(http.StatusOK, r.JSON(nil))
	}

	twitterClient := c.Value("twitterClient").(*TwitterClient)
	for _, twt := range actionableTweets {
		// https://eli.thegreenplace.net/2019/go-internals-capturing-loop-variables-in-closures/
		go func(twt tweet) {
			retry(func() error { return ingestTweet(twitterClient, &twt) })
		}(twt)
	}

	var ids []int64
	for _, tweet := range actionableTweets {
		ids = append(ids, tweet.ID)
	}
	response := map[string][]int64{
		"ingested_tweet_ids": ids,
	}
	return c.Render(http.StatusAccepted, r.JSON(response))
}

// validate that a tweet from our webhook is actionable
func validate(tweet *tweet) error {
	// tweet must be a reply
	if tweet.InReplyToStatusID == 0 {
		return errors.New("not a reply")
	}

	// tweet must not be a reply to lickerbot
	if tweet.InReplyToUserID == lickerbotUserID {
		return errors.New("lickerbot is not a bootlicker")
	}

	// tweet must mention lickerbot
	foundLickerbotMention := false
	for _, userMention := range tweet.Entities.UserMentions {
		if userMention.ID == lickerbotUserID {
			foundLickerbotMention = true
			break
		}
	}
	if !foundLickerbotMention {
		return errors.New("lickerbot not mentioned in tweet")
	}

	return nil
}

func ingestTweet(twitterClient *TwitterClient, tweet *tweet) error {
	log := logrus.WithField("tweetID", tweet.ID)
	log.Info("ingesting tweet")

	// create the bootlicker and lick in a transaction
	// we intentionally create these two in a transaction *before* sending our tweet
	// because our tweet will refer to the created bootlicker resource, e.g.
	// "lickerbot.com/@newuser". if the tweet is sent before those models are committed
	// there's a *chance* that the tweet will send and twitter will fetch a twitter card
	// for the page, but either due to error or lag the page won't actaully exist.
	bootlicker := &models.Bootlicker{}
	err := models.DB.Transaction(func(tx *pop.Connection) error {
		// fetch the bootlicker from our db. if they don't exist, create the record.
		// (also, i am very bad at nested error codeflow in go)
		err := tx.Where("lower(twitter_handle) = lower(?)", tweet.InReplyToScreenName).First(bootlicker)
		switch errors.Cause(err) {
		case nil:
		case sql.ErrNoRows:
			log.Info("bootlicker doesn't exist, creating")
			bootlicker = &models.Bootlicker{
				TwitterUserID: tweet.InReplyToUserID,
				TwitterHandle: tweet.InReplyToScreenName,
			}
			if err := tx.Create(bootlicker); err != nil {
				return err
			}
		default:
			return err
		}

		err = tx.Where("tweet_id = ?", tweet.InReplyToStatusID).First(&models.Lick{})
		switch errors.Cause(err) {
		case nil:
		case sql.ErrNoRows:
			lick := &models.Lick{
				TweetID:      tweet.InReplyToStatusID,
				TweetText:    tweet.Text,
				BootlickerID: bootlicker.ID,
			}
			if err := tx.Create(lick); err != nil {
				return err
			}
			if err := tx.Eager().Reload(bootlicker); err != nil {
				return err
			}
		default:
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	reply := buildReply(tweet, bootlicker)
	log.Infof("sending reply to twitter: %+v", reply)
	if err := twitterClient.TweetReply(reply); err != nil {
		return err
	}

	return nil
}

func buildReply(tweet *tweet, bootlicker *models.Bootlicker) *TweetReplyInput {
	times := "time"
	if len(bootlicker.Licks) > 1 {
		times = "times"
	}

	pledgeText := "Will you pledge to donate $5 per lick to organizations fighting police brutality?"
	if bootlicker.TotalPledged() > 0 {
		pledgeText = fmt.Sprintf(
			"In response, people have donated $%d to organizations fighting police brutality. That's $%d per lick!",
			bootlicker.TotalPledged(),
			int(math.Ceil(bootlicker.PledgedPerLick())),
		)
	}

	status := fmt.Sprintf(
		"@%s %s %s has licked the boot %d %s. %s https://lickerbot.com/@%s",
		tweet.User.ScreenName,
		lickReaction(tweet),
		bootlicker.TwitterHandle,
		len(bootlicker.Licks),
		times,
		pledgeText,
		bootlicker.TwitterHandle,
	)
	return &TweetReplyInput{
		Status:            status,
		InReplyToStatusID: tweet.ID,
	}
}

// returns a random lick reaction, based on the tweet's ID, for some testability.
func lickReaction(tweet *tweet) string {
	index := rand.New(rand.NewSource(tweet.InReplyToStatusID)).Intn(len(lickReactions))
	return lickReactions[index]
}

// retry a function if it errors.
func retry(fn func() error) error {
	log := logrus.WithField("fn", fn)
	var err error
	for attempt := 1; attempt < 4; attempt++ {
		if err = fn(); err == nil {
			return nil
		}
		if attempt < 4 {
			log.WithField("attempt", attempt).WithError(err).Warnf("error calling function. will sleep for %d seconds then retry.", attempt)

			sleepDuration := time.Duration(attempt) * time.Second
			time.Sleep(sleepDuration)
		}
	}
	log.WithError(err).Error("failed to successfully call function after retries")
	return err
}

// https://www.jokecamp.com/blog/examples-of-creating-base64-hashes-using-hmac-sha256-in-different-languages/#go
func hmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
