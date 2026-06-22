export const TIME_FORMAT = new Intl.DateTimeFormat("en", {
  month: "short",
  day: "numeric",
  year: "numeric",
  hour: "numeric",
  minute: "2-digit",
});

export const DATE_FORMAT = new Intl.DateTimeFormat("en", {
  month: "short",
  day: "numeric",
  year: "numeric",
});

export const formatTime = (date: number | string | Date): string =>
  TIME_FORMAT.format(new Date(date));

export const formatDate = (date: number | string | Date): string =>
  DATE_FORMAT.format(new Date(date));
