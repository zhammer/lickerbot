Feature: About Page
    I can see information about lickerbot on the about page.

    Scenario: I visit the About Page
        When I visit "/"
        Then I see the text "pledged to donate $20"
        And I see the text "$6.67 per lick"

    Scenario: I am a web crawler
        When I visit "/"
        Then the title is "Lickerbot"
        And the meta tag "og:title" equals "Lickerbot"
        And the meta tag "og:type" equals "website"
        And the meta tag "og:description" contains "Raise money for"
        And the meta tag "twitter:card" equals "summary"
        And the meta tag "twitter:site" equals "@lickerbot"
        And the meta tag "twitter:creator" equals "@zhammer_"
        And the meta tag "og:url" equals "https://lickerbot.com/"
        And the meta tag "og:image" contains a resource that exists
