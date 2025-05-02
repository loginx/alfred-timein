import tzLookup from 'tz-lookup';

/**
 * Get IANA timezone from coordinates.
 * Throws if not found.
 */
export function getTimezone(lat, lon) {
  if (typeof lat !== 'number' || typeof lon !== 'number') throw new Error('Latitude and longitude required');
  try {
    return tzLookup(lat, lon);
  } catch (err) {
    throw new Error('Timezone not found for coordinates');
  }
} 