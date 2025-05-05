import { existsSync, readFileSync, writeFileSync, mkdirSync, unlinkSync } from 'node:fs';
import { join } from 'node:path';
import { LRUCache } from 'lru-cache';

const CACHE_DIR = './.cache';
const CACHE_FILE = join(CACHE_DIR, 'cache.json');
const MAX_SIZE = 100;

// Ensure cache directory exists
mkdirSync(CACHE_DIR, { recursive: true });

function loadCache() {
  if (existsSync(CACHE_FILE)) {
    try {
      const data = JSON.parse(readFileSync(CACHE_FILE, 'utf8'));
      const lru = new LRUCache({ max: MAX_SIZE });
      if (Array.isArray(data.cache)) {
        for (const [key, value] of data.cache) {
          lru.set(key, value);
        }
      }
      return lru;
    } catch (e) {
      // Corrupt cache file, start fresh
    }
  }
  return new LRUCache({ max: MAX_SIZE });
}

function persistCache(cache) {
  const data = {
    max: cache.max,
    // Dump the cache as an array of [key, value] pairs
    cache: Array.from(cache.entries()),
  };
  writeFileSync(CACHE_FILE, JSON.stringify(data), 'utf8');
}

const cache = loadCache();

export async function getCache(key) {
  return cache.get(key);
}

export async function setCache(key, value) {
  cache.set(key, value);
  persistCache(cache);
}

export async function clearCache() {
  cache.clear();
  if (existsSync(CACHE_FILE)) unlinkSync(CACHE_FILE);
}

// For future: could add methods to clear or inspect the cache if needed.