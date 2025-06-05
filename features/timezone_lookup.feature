Feature: Timezone Lookup by City
  As a user who works with global teams
  I want to quickly find timezone information for any city
  So that I can coordinate meetings and deadlines across time zones

  Background:
    Given the timezone lookup service is available

  Scenario: Finding timezone for a major city
    Given I want to know the timezone for "Tokyo"
    When I request the timezone information
    Then I should get "Asia/Tokyo" as the timezone
    And the response should be fast

  Scenario: Finding timezone for a city with spaces in name
    Given I want to know the timezone for "New York"
    When I request the timezone information  
    Then I should get "America/New_York" as the timezone

  Scenario: Finding timezone for a landmark
    Given I want to know the timezone for "Eiffel Tower"
    When I request the timezone information
    Then I should get "Europe/Paris" as the timezone

  Scenario: Handling unknown locations gracefully
    Given I want to know the timezone for "NotARealCity123"
    When I request the timezone information
    Then I should receive a helpful error message
    And the error should mention "could not geocode"

  Scenario: Cache improves performance for repeated queries
    Given I have previously looked up "London"
    When I request the timezone for "London" again
    Then the response should be nearly instantaneous
    And I should get "Europe/London" as the timezone
    And the result should indicate it came from cache

  Scenario: Case insensitive city lookup
    Given I want to know the timezone for "BERLIN"
    When I request the timezone information
    Then I should get "Europe/Berlin" as the timezone

  Scenario: Whitespace handling
    Given I want to know the timezone for "  Paris  "
    When I request the timezone information
    Then I should get "Europe/Paris" as the timezone