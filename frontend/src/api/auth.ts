import { apiClient } from "./client";
import type { LoginRequest, LoginResponse } from "../types/models";

export const authApi = {
  async login(payload: LoginRequest) {
    const { data } = await apiClient.post<LoginResponse>("/api/auth/login", payload);
    return data;
  },
};
