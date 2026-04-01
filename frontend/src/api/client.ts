import axios from "axios";

import { getStoredToken } from "../utils/storage";

export const apiClient = axios.create({
  baseURL: "",
  headers: {
    "Content-Type": "application/json",
  },
});

apiClient.interceptors.request.use((config) => {
  const token = getStoredToken();

  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }

  return config;
});

export function extractApiError(error: unknown) {
  if (axios.isAxiosError(error)) {
    return (
      (error.response?.data as { error?: string } | undefined)?.error ??
      error.message ??
      "Ошибка запроса"
    );
  }

  if (error instanceof Error) {
    return error.message;
  }

  return "Неизвестная ошибка";
}
