package grifts

import (
	"fmt"
	"lickerbot/actions"

	"github.com/markbates/grift/grift"
)

var _ = grift.Namespace("twitter", func() {
	grift.Desc("register", "Register our webhook")
	grift.Add("register", func(c *grift.Context) error {
		twitter := actions.NewTwitterClient()
		webhookID, err := twitter.RegisterWebhook()
		if err != nil {
			return err
		}

		fmt.Println("Registered webhook:", webhookID)
		return nil
	})

	grift.Desc("unregister", "Unreigster a webhook")
	grift.Add("unregister", func(c *grift.Context) error {
		twitter := actions.NewTwitterClient()
		webhookID := c.Args[0]
		if err := twitter.UnregisterWebhook(webhookID); err != nil {
			return err
		}

		fmt.Println("Unregistered webhook:", webhookID)
		return nil
	})

	grift.Desc("subscribe", "Subscribes to account activity using registered webhook")
	grift.Add("subscribe", func(c *grift.Context) error {
		twitter := actions.NewTwitterClient()
		if err := twitter.SubscribeToAccountActivity(); err != nil {
			return err
		}

		fmt.Println("Successfully subscribed to account activity!")
		return nil
	})

	grift.Desc("tweet", "Fetch a tweet by its ID")
	grift.Add("tweet", func(c *grift.Context) error {
		twitter := actions.NewTwitterClient()
		tweet, err := twitter.FetchTweet(c.Args[0])
		if err != nil {
			return err
		}

		fmt.Println(tweet)
		return nil
	})
})
