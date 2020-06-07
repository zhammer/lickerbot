package actions

import (
	"lickerbot/models"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	totalLicks := 0
	totalLicks, _ = tx.Count(&models.Lick{})

	totalPledged := 0
	tx.RawQuery("SELECT SUM(amount) FROM pledged_donations").First(&totalPledged)

	pledgedPerLick := 0.0
	if totalLicks > 0 {
		pledgedPerLick = float64(totalPledged) / float64(totalLicks)
	}

	c.Set("totalLicks", totalLicks)
	c.Set("totalPledged", totalPledged)
	c.Set("pledgedPerLick", pledgedPerLick)
	return c.Render(http.StatusOK, r.HTML("index.html"))
}
