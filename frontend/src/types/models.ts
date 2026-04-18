export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  role: string;
}

export interface Warehouse {
  id: string;
  name: string;
  office_from_id: string;
  address: string;
  created_at: string;
}

export interface CreateWarehouseRequest {
  name: string;
  office_from_id: string;
  address: string;
}

export interface RouteItem {
  id: string;
  warehouse_id: string;
  route_id: string;
  name: string;
  created_at: string;
}

export interface CreateRouteRequest {
  route_id: string;
  name: string;
}

export interface Threshold {
  id: string;
  warehouse_id: string;
  route_id: string;
  value: number;
  updated_at: string;
}

export interface UpdateThresholdRequest {
  warehouse_id: string;
  route_id: string;
  value: number;
}

export interface HealthStatus {
  status: string;
}

export interface DataPoint {
  route_id: string;
  office_from_id: string;
  timestamp: string;
  status_1?: number | null;
  status_2?: number | null;
  status_3?: number | null;
  status_4?: number | null;
  status_5?: number | null;
  status_6?: number | null;
  status_7?: number | null;
  status_8?: number | null;
  target_2h?: number | null;
}

export interface DataIngestPayload {
  data_points: DataPoint[];
}

export interface DataIngestResponse {
  inserted: number;
}
