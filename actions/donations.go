package actions

import (
	"lickerbot/models"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/uuid"
)

type DonationRequest struct {
	Amount int `json:"amount"`
}

// DonationPledgeHandler shows information about a bootlicker.
func DonationPledgeHandler(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	donationPledgeRequest := &DonationRequest{}
	err := c.Bind(donationPledgeRequest)
	if err != nil {
		return err
	}

	bootlickerID, err := uuid.FromString(c.Param("bootlickerID"))
	if err != nil {
		return err
	}

	donationPledge := &models.PledgedDonation{
		Amount:       donationPledgeRequest.Amount,
		BootlickerID: bootlickerID,
	}
	err = tx.Create(donationPledge)
	if err != nil {
		return err
	}

	return c.Render(http.StatusCreated, r.JSON(nil))
}
