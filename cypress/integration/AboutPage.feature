Feature: About Page
    I can see information about lickerbot on the about page.

    Scenario: I visit the About Page
        When I visit "/"
        Then I see the text "pledged to donate $20"
        And I see the text "$6.67 per lick"
