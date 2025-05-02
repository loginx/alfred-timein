/**
 * Format the current time in a given IANA timezone.
 * Returns a string like 'Asia/Bangkok – Fri, May 2, 9:30 AM'.
 */
export function formatTime(timezone) {
  if (!timezone || typeof timezone !== 'string') throw new Error('Timezone required');
  const now = new Date();
  const options = {
    weekday: 'short',
    month: 'short',
    day: 'numeric',
    hour: 'numeric',
    minute: '2-digit',
    hour12: true,
  };
  const formatted = new Intl.DateTimeFormat('en-US', {
    ...options,
    timeZone: timezone,
  }).format(now);
  return `${timezone} – ${formatted}`;
} 