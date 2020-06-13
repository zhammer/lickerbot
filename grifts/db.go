package grifts

import (
	"lickerbot/models"

	"github.com/gobuffalo/pop/v5"
	"github.com/markbates/grift/grift"
)

var _ = grift.Namespace("db", func() {

	grift.Desc("seed", "Seeds a database")
	grift.Add("seed", func(c *grift.Context) error {
		return models.DB.Transaction(func(tx *pop.Connection) error {
			// clear tables
			if err := tx.RawQuery("TRUNCATE TABLE bootlickers CASCADE").Exec(); err != nil {
				return err
			}

			// bootlicker 1
			// could maybe use tx.Eager().Create() here to create nested models in one statement
			bootlicker1 := &models.Bootlicker{
				TwitterUserID: 929784070086672384,
				TwitterHandle: "MarcosDarkos",
			}
			if err := tx.Create(bootlicker1); err != nil {
				return err
			}

			lick1 := &models.Lick{
				BootlickerID: bootlicker1.ID,
				TweetID:      1268717186341916672,
				TweetText:    "Shouldn’t have tried to be a hero",
			}
			if err := tx.Create(lick1); err != nil {
				return err
			}

			lick2 := &models.Lick{
				BootlickerID: bootlicker1.ID,
				TweetID:      1268562799732629506,
				TweetText:    "U did nothing wrong, this is BS. There’s a reason the 1st amendment is first. Sucks u are forced to do this!",
			}
			if err := tx.Create(lick2); err != nil {
				return err
			}

			donation1 := &models.PledgedDonation{
				Amount:       15,
				BootlickerID: bootlicker1.ID,
			}
			if err := tx.Create(donation1); err != nil {
				return err
			}

			donation2 := &models.PledgedDonation{
				Amount:       5,
				BootlickerID: bootlicker1.ID,
			}
			if err := tx.Create(donation2); err != nil {
				return err
			}

			//bootlicker 2
			bootlicker2 := &models.Bootlicker{
				TwitterUserID: 1152084374,
				TwitterHandle: "OliverJWeber",
			}
			if err := tx.Create(bootlicker2); err != nil {
				return err
			}

			lick3 := &models.Lick{
				BootlickerID: bootlicker2.ID,
				TweetID:      1268556742251761666,
				TweetText:    "I sympathise with you when it comes to excess violence. But when you are restrained by a cop, you do *not* break away. That is resistance. The cop is not a man groping a woman, he is an officer of the law exercising a lawful extent of violence in restraining her movements.",
			}
			if err := tx.Create(lick3); err != nil {
				return err
			}

			return nil
		})
	})

})
