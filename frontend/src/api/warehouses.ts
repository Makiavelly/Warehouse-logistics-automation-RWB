import { apiClient } from "./client";
import type {
  CreateRouteRequest,
  CreateWarehouseRequest,
  RouteItem,
  Warehouse,
} from "../types/models";

export const warehousesApi = {
  async getAll() {
    const { data } = await apiClient.get<Warehouse[]>("/api/warehouses");
    return data;
  },
  async create(payload: CreateWarehouseRequest) {
    const { data } = await apiClient.post<Warehouse>("/api/warehouses", payload);
    return data;
  },
  async getRoutes(warehouseId: string) {
    const { data } = await apiClient.get<RouteItem[]>(`/api/warehouses/${warehouseId}/routes`);
    return data;
  },
  async createRoute(warehouseId: string, payload: CreateRouteRequest) {
    const { data } = await apiClient.post<RouteItem>(
      `/api/warehouses/${warehouseId}/routes`,
      payload,
    );
    return data;
  },
};
