package actions

import (
	"database/sql"
	"lickerbot/models"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"
)

// BootlickerHandler shows information about a bootlicker.
func BootlickerHandler(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	twitterHandle := c.Param("twitterHandle")

	bootlicker := models.Bootlicker{}
	err := tx.Eager().Where("lower(twitter_handle) = lower(?)", twitterHandle).First(&bootlicker)
	if err != nil && errors.Cause(err) == sql.ErrNoRows {
		c.Set("twitterHandle", twitterHandle)
		return c.Render(http.StatusNotFound, r.HTML("bootlicker_not_found.html"))
	}
	if err != nil {
		return err
	}

	c.Set("bootlicker", bootlicker)
	return c.Render(http.StatusOK, r.HTML("bootlicker.html"))
}
