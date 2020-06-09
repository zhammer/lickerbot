Feature: Bootlicker Page
    I can see information about individual bootlickers on the bootlicker page.

    Scenario: I visit the Bootlicker Page for a bootlicker with donations
        When I visit "/@MarcosDarkos"
        Then I see the text "2 times"
        And I see the text "people have donated $20"
        And I see the text "$10 per lick"
        And I see 2 embedded tweets

    Scenario: I visit the Bootlicker Page for a bootlicker without donations
        When I visit "/@OliverJWeber"
        Then I see the text "1 time"
        And I see the text "Will you pledge to donate $5 per lick?"
        And I see 1 embedded tweets
        And I see the text "I sympathise with you"

    Scenario: I visit the Bootlicker Page for a non-existing bootlicker
        When I visit "/@zhammer_" expecting a non-200 status
        Then I see the text "@zhammer_ has not yet been caught bootlicking"
        And I see 0 embedded tweets

    Scenario: I navigate to the About Page
        When I click the link "Lickerbot"
        Then I am on "/"

    Scenario: I pledge to donate
        When I visit "/@marcosdarkos"
        And I click the link "join them"
        Then I am on the hash "#donate"
        And I see the text "Pledge to donate $10"

        When I select the amount "$20"
        Then I see the text "Pledge to donate $40"

        When I click the button "Pledge to donate $40"
        Then I see the text "Thank you for pledging to donate $40"

        When I refresh the page
        Then I see the text "people have donated $60"

    Scenario: I am a webcrawler
        When I visit "/@MarcosDarkos"
        Then the title is "Lickerbot - MarcosDarkos"
        And the meta tag "og:title" equals "Lickerbot - MarcosDarkos"
        And the meta tag "og:url" equals "https://lickerbot.com/@MarcosDarkos/"
