import { apiClient } from "./client";
import type { DataIngestPayload, DataIngestResponse } from "../types/models";

export const dataIngestApi = {
  async ingest(payload: DataIngestPayload, apiKey: string) {
    const { data } = await apiClient.post<DataIngestResponse>("/api/data/ingest", payload, {
      headers: {
        "X-API-Key": apiKey,
      },
    });

    return data;
  },
};
