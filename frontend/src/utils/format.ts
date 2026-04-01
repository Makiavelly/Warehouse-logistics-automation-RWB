import dayjs from "dayjs";

export function formatDateTime(value?: string | Date | null) {
  if (!value) {
    return "Нет данных";
  }

  return dayjs(value).format("DD.MM.YYYY HH:mm");
}

export function formatDate(value?: string | Date | null) {
  if (!value) {
    return "Нет данных";
  }

  return dayjs(value).format("DD.MM.YYYY");
}

export function formatNumber(value?: number | null, digits = 0) {
  if (value === null || value === undefined || Number.isNaN(value)) {
    return "0";
  }

  return new Intl.NumberFormat("ru-RU", {
    maximumFractionDigits: digits,
    minimumFractionDigits: digits,
  }).format(value);
}

export function truncateMiddle(value: string, size = 18) {
  if (value.length <= size) {
    return value;
  }

  const edge = Math.max(4, Math.floor((size - 1) / 2));
  return `${value.slice(0, edge)}…${value.slice(-edge)}`;
}
