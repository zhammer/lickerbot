package actions

import (
	"math"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	forcessl "github.com/gobuffalo/mw-forcessl"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/unrolled/secure"

	"lickerbot/models"

	"github.com/gobuffalo/buffalo-pop/v2/pop/popmw"
	csrf "github.com/gobuffalo/mw-csrf"
	i18n "github.com/gobuffalo/mw-i18n"
	"github.com/gobuffalo/packr/v2"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var T *i18n.Translator

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
//
// Routing, middleware, groups, etc... are declared TOP -> DOWN.
// This means if you add a middleware to `app` *after* declaring a
// group, that group will NOT have that new middleware. The same
// is true of resource declarations as well.
//
// It also means that routes are checked in the order they are declared.
// `ServeFiles` is a CATCH-ALL route, so it should always be
// placed last in the route declarations, as it will prevent routes
// declared after it to never be called.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_lickerbot_session",
		})

		// Automatically redirect to SSL
		app.Use(forceSSL())

		// Log request parameters (filters apply).
		app.Use(paramlogger.ParameterLogger)

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(csrf.New)

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.Connection)
		// Remove to disable this.
		app.Use(popmw.Transaction(models.DB))

		// Setup and use translations:
		app.Use(translations())

		app.Use(twitterMw())

		app.Use(templateGlobals)

		app.GET("/", HomeHandler)

		app.GET("/@{twitterHandle}", BootlickerHandler)

		app.POST("/{bootlickerID}/donations", DonationPledgeHandler)

		app.GET("/webhook/twitter", TwitterSecurityCheck)

		app.POST("/webhook/twitter", TwitterWebhook)

		app.ServeFiles("/", assetsBox) // serve files from the public directory
	}

	return app
}

// translations will load locale files, set up the translator `actions.T`,
// and will return a middleware to use to load the correct locale for each
// request.
// for more information: https://gobuffalo.io/en/docs/localization
func translations() buffalo.MiddlewareFunc {
	var err error
	if T, err = i18n.New(packr.New("app:locales", "../locales"), "en-US"); err != nil {
		app.Stop(err)
	}
	return T.Middleware()
}

// forceSSL will return a middleware that will redirect an incoming request
// if it is not HTTPS. "http://example.com" => "https://example.com".
// This middleware does **not** enable SSL. for your application. To do that
// we recommend using a proxy: https://gobuffalo.io/en/docs/proxy
// for more information: https://github.com/unrolled/secure/
func forceSSL() buffalo.MiddlewareFunc {
	return forcessl.Middleware(secure.Options{
		SSLRedirect:     ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}

type link struct {
	Name string
	URL  string
}

// templateGlobals adds some global values to all templates.
func templateGlobals(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		c.Set("donationLinks", []link{
			{
				Name: "Black Trans Protestors Emergency Fund",
				URL:  "https://www.instagram.com/p/CA8GE-HDbxa/",
			},
			{
				Name: "Justice for Breonna Taylor",
				URL:  "https://www.gofundme.com/f/9v4q2-justice-for-breonna-taylor",
			},
			{
				Name: "The Bail Project",
				URL:  "https://secure.givelively.org/donate/the-bail-project",
			},
			{
				Name: "More",
				URL:  "https://blacklivesmatters.carrd.co/#donate",
			},
		})
		c.Set("truncateFloat", truncateFloat)
		c.Set("url", c.Request().URL)
		return next(c)
	}
}

// 10.25897237 -> 10.26
func truncateFloat(input float64) float64 {
	return math.Ceil(input*100) / 100
}

// twitter middleware, attaches a twitter client to request contexts
func twitterMw() func(buffalo.Handler) buffalo.Handler {
	client := NewTwitterClient()
	return func(next buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {
			c.Set("twitterClient", client)
			return next(c)
		}
	}
}
