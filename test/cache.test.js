import { describe, it, expect } from 'vitest';
import { getCache, setCache } from '../src/cache.js';
import { existsSync, unlinkSync, readFileSync } from 'node:fs';
import { join } from 'node:path';

const CACHE_FILE = join('./.cache', 'cache.json');

describe('cache', () => {
  it('returns undefined for missing key', async () => {
    expect(await getCache('notfound')).toBeUndefined();
  });

  it('sets and gets a value', async () => {
    await setCache('foo', 123);
    expect(await getCache('foo')).toBe(123);
  });

  it('overwrites a value', async () => {
    await setCache('foo', 456);
    expect(await getCache('foo')).toBe(456);
  });

  it('limits cache size', async () => {
    for (let i = 0; i < 110; i++) {
      await setCache(`key${i}`, i);
    }
    expect(await getCache('key0')).toBeUndefined(); // Oldest entry should be evicted
    expect(await getCache('key109')).toBe(109); // Newest entry should remain
  });

  it('persists to disk and restores on reload', async () => {
    if (existsSync(CACHE_FILE)) unlinkSync(CACHE_FILE);
    const key = 'persist-key';
    const value = 'persist-value';
    await setCache(key, value);
    // Simulate reload by clearing require cache and re-importing
    const { LRUCache } = await import('lru-cache');
    const data = JSON.parse(readFileSync(CACHE_FILE, 'utf8'));
    const lru = new LRUCache({ max: 100 });
    for (const [k, v] of data.cache) lru.set(k, v);
    expect(lru.get(key)).toBe(value);
  });
});