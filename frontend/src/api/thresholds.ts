import { apiClient } from "./client";
import type { Threshold, UpdateThresholdRequest } from "../types/models";

export const thresholdsApi = {
  async getAll(filters?: { warehouseId?: string; routeId?: string }) {
    const params = new URLSearchParams();

    if (filters?.warehouseId) {
      params.set("warehouse_id", filters.warehouseId);
    }

    if (filters?.routeId) {
      params.set("route_id", filters.routeId);
    }

    const query = params.size > 0 ? `?${params.toString()}` : "";
    const { data } = await apiClient.get<Threshold[]>(`/api/thresholds${query}`);
    return data;
  },
  async update(payload: UpdateThresholdRequest) {
    const { data } = await apiClient.put<Threshold>("/api/thresholds", payload);
    return data;
  },
};
