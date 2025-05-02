import { describe, it, expect, beforeAll, afterAll } from 'vitest';
import { formatTime } from '../src/formatTime.js';

const fixedDate = new Date('2024-05-02T02:30:00Z');
let RealDate;

describe('formatTime', () => {
  beforeAll(() => {
    RealDate = global.Date;
    class MockDate extends Date {
      constructor(...args) {
        if (args.length === 0) {
          super(fixedDate.getTime());
        } else {
          super(...args);
        }
      }
    }
    global.Date = MockDate;
  });

  afterAll(() => {
    global.Date = RealDate;
  });

  it('formats time for a valid timezone', () => {
    const result = formatTime('Asia/Bangkok');
    expect(result).toMatch(/^Asia\/Bangkok – \w{3}, \w{3} \d, \d{1,2}:\d{2} (AM|PM)$/);
  });

  it('throws for missing timezone', () => {
    expect(() => formatTime()).toThrow('Timezone required');
  });

  it('throws for empty string', () => {
    expect(() => formatTime('')).toThrow('Timezone required');
  });
}); 