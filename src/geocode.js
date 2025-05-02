import NodeGeocoder from 'node-geocoder';

const geocoder = NodeGeocoder({
  provider: 'openstreetmap',
  // No API key needed
});

/**
 * Resolve a city name to { lat, lon }.
 * Throws if not found.
 */
export async function getCoordinates(city) {
  if (!city || typeof city !== 'string') throw new Error('City name required');
  const results = await geocoder.geocode(city);
  if (!results.length) throw new Error(`City not found: ${city}`);
  // Pick the first result
  const { latitude: lat, longitude: lon } = results[0];
  if (typeof lat !== 'number' || typeof lon !== 'number') throw new Error(`Invalid coordinates for: ${city}`);
  return { lat, lon };
} 