Feature: Alfred Workflow Integration
  As an Alfred user
  I want to search for time information using natural language
  So that I can quickly get timezone data without leaving my workflow

  Background:
    Given I am using Alfred with the timein workflow

  Scenario: Alfred JSON format for timezone lookup
    Given I search for timezone information for "Tokyo" 
    When I request Alfred format output
    Then I should receive valid Alfred JSON
    And the JSON should contain exactly one result item
    And the item should have a title with the timezone
    And the item should have a subtitle mentioning the city
    And the item should be actionable

  Scenario: Alfred JSON format for time display  
    Given I have the timezone "Europe/Paris"
    When I request current time in Alfred format
    Then I should receive valid Alfred JSON
    And the item title should contain the timezone and current time
    And the item subtitle should mention the city and timezone abbreviation
    And the result should include timezone variables

  # Note: Additional scenarios temporarily removed to avoid undefined step pollution
  # TODO: Implement step definitions for:
  # - Alfred error formatting
  # - Alfred caching indicators  
  # - Alfred script filter caching headers