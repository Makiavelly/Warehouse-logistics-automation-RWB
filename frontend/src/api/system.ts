import { apiClient } from "./client";
import type { HealthStatus } from "../types/models";

export const systemApi = {
  async health() {
    const { data } = await apiClient.get<HealthStatus>("/health");
    return data;
  },
};
