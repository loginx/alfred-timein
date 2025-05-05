// Suppress warnings globally for terminal execution
process.env.NODE_NO_WARNINGS = '1';

// Filter out warnings from the output
const originalConsoleWarn = console.warn;
console.warn = (...args) => {
  if (!args[0]?.includes('DeprecationWarning')) {
    originalConsoleWarn(...args);
  }
};

import alfy from 'alfy';
import { getCache, setCache } from './src/cache.js';
import { getCoordinates } from './src/geocode.js';
import { getTimezone } from './src/timezone.js';
import { formatTime } from './src/formatTime.js';

const input = (alfy.input || '').trim();

// Ensure only the latest input is processed
let lastInput = '';

async function main() {
  if (input === lastInput) return; // Ignore duplicate inputs
  lastInput = input;

  if (!input) {
    alfy.output([
      {
        title: 'Enter a city name',
        subtitle: 'Example: timein Bangkok',
        valid: false,
        icon: { path: alfy.icon.info },
      },
    ]);
    return;
  }

  // Check cache first
  let timezone = await getCache(input.toLowerCase());
  if (!timezone) {
    let coords;
    try {
      coords = await getCoordinates(input);
    } catch (err) {
      alfy.output([
        {
          title: 'City not found',
          subtitle: err.message,
          valid: false,
          icon: { path: alfy.icon.warning },
        },
      ]);
      return;
    }
    try {
      timezone = getTimezone(coords.lat, coords.lon);
      await setCache(input.toLowerCase(), timezone);
    } catch (err) {
      alfy.output([
        {
          title: 'Timezone not found',
          subtitle: err.message,
          valid: false,
          icon: { path: alfy.icon.warning },
        },
      ]);
      return;
    }
  }

  let timeString;
  try {
    timeString = formatTime(timezone);
  } catch (err) {
    alfy.output([
      {
        title: 'Could not format time',
        subtitle: err.message,
        valid: false,
        icon: { path: alfy.icon.error },
      },
    ]);
    return;
  }

  alfy.output([
    {
      title: timeString,
      subtitle: `Current time in ${input}`,
      arg: timeString,
      icon: { path: alfy.icon.get('Clock') },
    },
  ]);
}

main();

// UX improvement: Add fuzzy matching for city names (e.g., suggest close matches if not found)