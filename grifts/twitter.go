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
})
