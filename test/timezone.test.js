import { describe, it, expect, vi } from 'vitest';
import * as timezoneModule from '../src/timezone.js';

vi.mock('tz-lookup', () => ({
  default: (lat, lon) => {
    if (lat === 13.7563 && lon === 100.5018) return 'Asia/Bangkok';
    throw new Error('Not found');
  },
}));

describe('getTimezone', () => {
  it('returns timezone for valid coordinates', () => {
    expect(timezoneModule.getTimezone(13.7563, 100.5018)).toBe('Asia/Bangkok');
  });

  it('throws for invalid coordinates', () => {
    expect(() => timezoneModule.getTimezone(0, 0)).toThrow('Timezone not found');
  });

  it('throws for non-number input', () => {
    expect(() => timezoneModule.getTimezone('foo', 'bar')).toThrow('Latitude and longitude required');
  });
}); 