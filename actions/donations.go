package actions

import (
	"lickerbot/models"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/uuid"
)

type DonationRequest struct {
	Amount int `json:"amount"`
}

// DonationPledgeHandler shows information about a bootlicker.
func DonationPledgeHandler(c buffalo.Context) error {
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
	err = models.DB.Create(donationPledge)
	if err != nil {
		return err
	}

	return c.Render(http.StatusCreated, r.JSON(nil))
}
