import { describe, it, expect } from 'vitest';
import { getCache, setCache } from '../src/cache.js';

describe('cache', () => {
  it('returns undefined for missing key', () => {
    expect(getCache('notfound')).toBeUndefined();
  });

  it('sets and gets a value', () => {
    setCache('foo', 123);
    expect(getCache('foo')).toBe(123);
  });

  it('overwrites a value', () => {
    setCache('foo', 456);
    expect(getCache('foo')).toBe(456);
  });
}); 