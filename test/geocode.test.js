import { describe, it, expect, vi } from 'vitest';
import * as geocodeModule from '../src/geocode.js';

// Mock node-geocoder
vi.mock('node-geocoder', () => ({
  default: () => ({
    geocode: async (city) => {
      if (city === 'Bangkok') return [{ latitude: 13.7563, longitude: 100.5018 }];
      if (city === 'Nowhereville') return [];
      return [];
    },
  }),
}));

describe('getCoordinates', () => {
  it('returns coordinates for a valid city', async () => {
    const coords = await geocodeModule.getCoordinates('Bangkok');
    expect(coords).toEqual({ lat: 13.7563, lon: 100.5018 });
  });

  it('throws for an invalid city', async () => {
    await expect(geocodeModule.getCoordinates('Nowhereville')).rejects.toThrow('City not found');
  });

  it('throws for empty input', async () => {
    await expect(geocodeModule.getCoordinates('')).rejects.toThrow('City name required');
  });
});