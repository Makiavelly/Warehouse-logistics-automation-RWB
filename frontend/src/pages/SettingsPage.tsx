import CloudUploadRounded from "@mui/icons-material/CloudUploadRounded";
import KeyRounded from "@mui/icons-material/KeyRounded";
import MemoryRounded from "@mui/icons-material/MemoryRounded";
import { Alert, Box, Button, Stack, TextField, Typography } from "@mui/material";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useMemo, useState } from "react";

import { dataIngestApi } from "../api/dataIngest";
import { extractApiError } from "../api/client";
import { systemApi } from "../api/system";
import { useSession } from "../app/session";
import { GlassPanel } from "../components/GlassPanel";
import { KpiCard } from "../components/KpiCard";

const samplePayload = `{
  "data_points": [
    {
      "route_id": "route-101",
      "office_from_id": "WH-001",
      "timestamp": "2026-04-01T10:00:00Z",
      "status_1": 24,
      "status_2": 17,
      "target_2h": 56
    }
  ]
}`;

export function SettingsPage() {
  const { apiKey, setApiKey } = useSession();
  const [apiKeyDraft, setApiKeyDraft] = useState(apiKey);
  const [payload, setPayload] = useState(samplePayload);
  const [resultText, setResultText] = useState("");
  const [errorText, setErrorText] = useState("");

  const healthQuery = useQuery({
    queryKey: ["settings-health"],
    queryFn: systemApi.health,
  });

  const parsedPreview = useMemo(() => {
    try {
      const value = JSON.parse(payload) as unknown;

      if (Array.isArray(value)) {
        return value.length;
      }

      if (value && typeof value === "object" && Array.isArray((value as { data_points?: unknown[] }).data_points)) {
        return (value as { data_points: unknown[] }).data_points.length;
      }

      return 0;
    } catch {
      return 0;
    }
  }, [payload]);

  const ingestMutation = useMutation({
    mutationFn: async () => {
      const raw = JSON.parse(payload) as unknown;
      const normalized =
        Array.isArray(raw) ? { data_points: raw } : (raw as { data_points: unknown[] });

      return dataIngestApi.ingest(normalized, apiKeyDraft);
    },
    onSuccess: (response) => {
      setErrorText("");
      setResultText(`Успешно загружено ${response.inserted} записей.`);
      setApiKey(apiKeyDraft);
    },
    onError: (error) => {
      setResultText("");
      setErrorText(extractApiError(error));
    },
  });

  return (
    <Stack spacing={3}>
      <Box
        sx={{
          display: "grid",
          gap: 2,
          gridTemplateColumns: { xs: "1fr", md: "repeat(3, 1fr)" },
        }}
      >
        <KpiCard
          label="Health"
          value={healthQuery.data?.status === "ok" ? "Online" : "Offline"}
          caption="Проверка живости backend"
          icon={<MemoryRounded />}
        />
        <KpiCard
          label="Preview datapoints"
          value={String(parsedPreview)}
          caption="Столько элементов сейчас распознано в JSON"
          icon={<CloudUploadRounded />}
        />
        <KpiCard
          label="API key profile"
          value={apiKeyDraft ? "Configured" : "Missing"}
          caption="X-API-Key для POST /api/data/ingest"
          icon={<KeyRounded />}
        />
      </Box>

      <Box
        sx={{
          display: "grid",
          gap: 2,
          gridTemplateColumns: "1fr",
        }}
      >
        <GlassPanel>
          <Stack spacing={2}>
            <Typography variant="h5">Загрузка данных</Typography>
            <Typography color="text.secondary">
              Endpoint принимает JSON-объект с полем `data_points` и требует заголовок
              `X-API-Key`.
            </Typography>
            {resultText ? <Alert severity="success">{resultText}</Alert> : null}
            {errorText ? <Alert severity="error">{errorText}</Alert> : null}

            <TextField
              label="API key"
              value={apiKeyDraft}
              onChange={(event) => setApiKeyDraft(event.target.value)}
            />
            <TextField
              label="JSON payload"
              multiline
              minRows={14}
              value={payload}
              onChange={(event) => setPayload(event.target.value)}
            />
            <Stack direction={{ xs: "column", sm: "row" }} spacing={1.25}>
              <Button variant="contained" onClick={() => ingestMutation.mutate()} disabled={ingestMutation.isPending}>
                {ingestMutation.isPending ? "Загружаем..." : "Отправить в ingest"}
              </Button>
              <Button variant="outlined" onClick={() => setPayload(samplePayload)}>
                Вернуть sample
              </Button>
            </Stack>
          </Stack>
        </GlassPanel>
      </Box>
    </Stack>
  );
}
