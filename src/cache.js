import QuickLRU from 'quick-lru';

const lru = new QuickLRU({ maxSize: 100 });

export function getCache(key) {
  return lru.get(key);
}

export function setCache(key, value) {
  lru.set(key, value);
}

// For future: could persist to disk if needed, but in-memory is fast and simple. 