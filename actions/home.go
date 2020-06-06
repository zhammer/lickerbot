package actions

import (
	"lickerbot/models"
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	totalLicks := 0
	totalLicks, _ = models.DB.Count(&models.Lick{})

	totalPledged := 0
	models.DB.RawQuery("SELECT SUM(amount) FROM pledged_donations").First(&totalPledged)

	pledgedPerLick := 0.0
	if totalLicks > 0 {
		pledgedPerLick = float64(totalPledged) / float64(totalLicks)
	}

	c.Set("totalLicks", totalLicks)
	c.Set("totalPledged", totalPledged)
	c.Set("pledgedPerLick", pledgedPerLick)
	return c.Render(http.StatusOK, r.HTML("index.html"))
}
