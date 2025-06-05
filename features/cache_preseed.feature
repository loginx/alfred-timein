Feature: Cache Pre-seeding
  As a user of the timezone lookup tool
  I want major world capitals to be pre-cached
  So that I get fast responses for common cities

  Background:
    Given the cache has been pre-seeded with capital cities

  Scenario: Pre-seeded cache hit for major capital
    When I look up the timezone for "London"
    Then the response should be fast
    And the result should contain "Europe/London"
    And the cache should indicate a hit

  Scenario: Pre-seeded cache hit for Asian capital
    When I look up the timezone for "Tokyo"
    Then the response should be fast
    And the result should contain "Asia/Tokyo"
    And the cache should indicate a hit

  Scenario: Pre-seeded entries have long TTL
    Given the cache was pre-seeded 30 days ago
    When I look up the timezone for "Paris"
    Then the result should contain "Europe/Paris"
    And the cache should indicate a hit

  Scenario: User cache entries override pre-seeded entries
    Given the cache has a user-created entry for "London"
    When I look up the timezone for "London"
    Then the result should use the user-created entry
    And not the pre-seeded entry