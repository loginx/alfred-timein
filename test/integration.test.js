import { describe, it, expect } from 'vitest';
import { getCache, setCache, clearCache } from '../src/cache.js';
import { getCoordinates } from '../src/geocode.js';
import { getTimezone } from '../src/timezone.js';
import { unlinkSync, existsSync } from 'node:fs';
import { join } from 'node:path';

const CACHE_FILE = join('./.cache', 'cache.json');

describe('integration: cache and timezone', () => {
  it('uses the cache across requests for city->timezone', async () => {
    const city = 'Bangkok';
    const cityKey = city.toLowerCase();

    // Ensure cache is empty
    await clearCache();
    expect(await getCache(cityKey)).toBeUndefined();

    // First request: resolve coordinates, then timezone, then cache
    const coords = await getCoordinates(city);
    const timezone = getTimezone(coords.lat, coords.lon);
    await setCache(cityKey, timezone);
    expect(await getCache(cityKey)).toBe(timezone);

    // Second request: should use the cached timezone
    const cachedTimezone = await getCache(cityKey);
    expect(cachedTimezone).toBe(timezone);
  });
});
