import Inventory2Rounded from "@mui/icons-material/Inventory2Rounded";
import LanRounded from "@mui/icons-material/LanRounded";
import RouteRounded from "@mui/icons-material/RouteRounded";
import TuneRounded from "@mui/icons-material/TuneRounded";
import {
  Alert,
  Box,
  LinearProgress,
  Stack,
  Typography,
} from "@mui/material";
import { useQuery } from "@tanstack/react-query";

import { systemApi } from "../api/system";
import { thresholdsApi } from "../api/thresholds";
import { warehousesApi } from "../api/warehouses";
import { GlassPanel } from "../components/GlassPanel";
import { KpiCard } from "../components/KpiCard";
import { LoadingBlock } from "../components/LoadingBlock";
import { formatNumber } from "../utils/format";

export function DashboardPage() {
  const healthQuery = useQuery({
    queryKey: ["dashboard-health"],
    queryFn: systemApi.health,
  });

  const warehousesQuery = useQuery({
    queryKey: ["dashboard-warehouses"],
    queryFn: warehousesApi.getAll,
  });

  const thresholdsQuery = useQuery({
    queryKey: ["dashboard-thresholds"],
    queryFn: () => thresholdsApi.getAll(),
  });

  const totalRoutesQuery = useQuery({
    queryKey: ["dashboard-route-count", warehousesQuery.data?.map((item) => item.id).join(",") ?? "none"],
    enabled: Boolean(warehousesQuery.data?.length),
    queryFn: async () => {
      const routes = await Promise.all(
        (warehousesQuery.data ?? []).map((warehouse) => warehousesApi.getRoutes(warehouse.id)),
      );

      return routes.reduce((sum, current) => sum + current.length, 0);
    },
  });

  const isLoading =
    healthQuery.isLoading || warehousesQuery.isLoading || thresholdsQuery.isLoading;

  const cards = [
    {
      label: "Backend статус",
      value: healthQuery.data?.status === "ok" ? "Online" : "Offline",
      caption: "Проверка через GET /health",
      icon: <LanRounded />,
    },
    {
      label: "Склады в системе",
      value: formatNumber(warehousesQuery.data?.length ?? 0),
      caption: "Текущий каталог активных складов",
      icon: <Inventory2Rounded />,
    },
    {
      label: "Маршруты в справочнике",
      value: formatNumber(totalRoutesQuery.data ?? 0),
      caption: "Сумма по всем складам",
      icon: <RouteRounded />,
    },
    {
      label: "Настроенные пороги",
      value: formatNumber(thresholdsQuery.data?.length ?? 0),
      caption: "GET /api/thresholds",
      icon: <TuneRounded />,
    },
  ];

  return (
    <Stack spacing={3}>
      {healthQuery.isError ? (
        <Alert severity="error">
          Backend сейчас недоступен. Каркас страниц продолжает работать, но живые блоки требуют
          запущенный `logistics-service`.
        </Alert>
      ) : null}

      <Box
        sx={{
          display: "grid",
          gap: 2,
          gridTemplateColumns: { xs: "1fr", md: "repeat(2, 1fr)", xl: "repeat(4, 1fr)" },
        }}
      >
        {cards.map((card) => (
          <KpiCard key={card.label} {...card} />
        ))}
      </Box>

      <Box
        sx={{
          display: "grid",
          gap: 2,
          gridTemplateColumns: "1fr",
        }}
      >
        <GlassPanel>
          <Stack spacing={2.5}>
            <Typography variant="h5">Покрытие складов порогами</Typography>

            {isLoading ? (
              <LoadingBlock />
            ) : (
              <Stack spacing={2}>
                {(warehousesQuery.data ?? []).slice(0, 6).map((warehouse) => {
                  const thresholdsForWarehouse =
                    thresholdsQuery.data?.filter((item) => item.warehouse_id === warehouse.id).length ?? 0;
                  const totalRoutes = totalRoutesQuery.data ? Math.max(totalRoutesQuery.data, 1) : 1;
                  const progress = Math.min(100, Math.round((thresholdsForWarehouse / totalRoutes) * 100));

                  return (
                    <Box key={warehouse.id}>
                      <Stack direction="row" justifyContent="space-between" sx={{ mb: 0.75 }}>
                        <Typography variant="body1">{warehouse.name}</Typography>
                        <Typography color="text.secondary" variant="body2">
                          {warehouse.office_from_id}
                        </Typography>
                      </Stack>
                      <LinearProgress variant="determinate" value={progress} sx={{ height: 10, borderRadius: 999 }} />
                    </Box>
                  );
                })}
              </Stack>
            )}
          </Stack>
        </GlassPanel>
      </Box>
    </Stack>
  );
}
