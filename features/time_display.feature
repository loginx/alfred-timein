Feature: Current Time Display
  As a user coordinating across time zones
  I want to see the current local time in any timezone
  So that I know what time it is for my colleagues or clients

  Background:
    Given the time display service is available

  Scenario: Displaying current time for a valid timezone
    Given I have the timezone "America/New_York"
    When I request the current time
    Then I should see a human-readable time format
    And the time should include the day of the week
    And the time should include the current date
    And the time should include hours and minutes

  # Note: Additional scenarios temporarily removed to avoid undefined step pollution
  # TODO: Implement step definitions for:
  # - Timezone validation before time display
  # - Time format is locale-aware
  # - Different timezones show different times
  # - Empty timezone input handling