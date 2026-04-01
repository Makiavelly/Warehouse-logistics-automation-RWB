import AddRounded from "@mui/icons-material/AddRounded";
import TuneRounded from "@mui/icons-material/TuneRounded";
import WarehouseRounded from "@mui/icons-material/WarehouseRounded";
import RouteRounded from "@mui/icons-material/RouteRounded";
import LayersRounded from "@mui/icons-material/LayersRounded";
import {
  Alert,
  Box,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  MenuItem,
  Stack,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TextField,
  Typography,
} from "@mui/material";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useMemo, useState } from "react";
import { Link as RouterLink, useParams } from "react-router-dom";

import { extractApiError } from "../api/client";
import { thresholdsApi } from "../api/thresholds";
import { warehousesApi } from "../api/warehouses";
import { EmptyState } from "../components/EmptyState";
import { GlassPanel } from "../components/GlassPanel";
import { KpiCard } from "../components/KpiCard";
import { LoadingBlock } from "../components/LoadingBlock";
import { formatDateTime, formatNumber } from "../utils/format";

export function WarehouseDetailsPage() {
  const params = useParams();
  const warehouseId = params.warehouseId ?? "";
  const queryClient = useQueryClient();
  const [routeDialogOpen, setRouteDialogOpen] = useState(false);
  const [thresholdError, setThresholdError] = useState("");
  const [routeError, setRouteError] = useState("");
  const [routeId, setRouteId] = useState("");
  const [routeName, setRouteName] = useState("");
  const [selectedRouteId, setSelectedRouteId] = useState("");
  const [thresholdValue, setThresholdValue] = useState("10");

  const warehousesQuery = useQuery({
    queryKey: ["warehouses"],
    queryFn: warehousesApi.getAll,
  });

  const routesQuery = useQuery({
    queryKey: ["routes", warehouseId],
    queryFn: () => warehousesApi.getRoutes(warehouseId),
    enabled: Boolean(warehouseId),
  });

  const thresholdsQuery = useQuery({
    queryKey: ["thresholds", warehouseId],
    queryFn: () => thresholdsApi.getAll({ warehouseId }),
    enabled: Boolean(warehouseId),
  });

  const warehouse = useMemo(
    () => warehousesQuery.data?.find((item) => item.id === warehouseId),
    [warehouseId, warehousesQuery.data],
  );

  const createRouteMutation = useMutation({
    mutationFn: (payload: { route_id: string; name: string }) =>
      warehousesApi.createRoute(warehouseId, payload),
    onSuccess: async () => {
      setRouteDialogOpen(false);
      setRouteId("");
      setRouteName("");
      await queryClient.invalidateQueries({ queryKey: ["routes", warehouseId] });
    },
    onError: (error) => setRouteError(extractApiError(error)),
  });

  const updateThresholdMutation = useMutation({
    mutationFn: thresholdsApi.update,
    onSuccess: async () => {
      setThresholdError("");
      await queryClient.invalidateQueries({ queryKey: ["thresholds", warehouseId] });
    },
    onError: (error) => setThresholdError(extractApiError(error)),
  });

  if (warehousesQuery.isLoading || routesQuery.isLoading || thresholdsQuery.isLoading) {
    return <LoadingBlock />;
  }

  if (!warehouse) {
    return (
      <EmptyState
        title="Склад не найден"
        description="Похоже, карточка склада была удалена или вы открыли неверный URL."
      />
    );
  }

  const thresholdsMap = new Map((thresholdsQuery.data ?? []).map((item) => [item.route_id, item]));

  return (
    <Stack spacing={3}>
      <Stack direction={{ xs: "column", md: "row" }} spacing={1.25}>
        <Button component={RouterLink} to="/warehouses" variant="outlined">
          Назад к списку
        </Button>
        <Button startIcon={<AddRounded />} variant="contained" onClick={() => setRouteDialogOpen(true)}>
          Добавить маршрут
        </Button>
      </Stack>

      <Box
        sx={{
          display: "grid",
          gap: 2,
          gridTemplateColumns: { xs: "1fr", md: "repeat(2, 1fr)", xl: "repeat(4, 1fr)" },
        }}
      >
        <KpiCard
          label="Маршрутов"
          value={formatNumber(routesQuery.data?.length ?? 0)}
          caption="GET /api/warehouses/{id}/routes"
          icon={<RouteRounded />}
        />
        <KpiCard
          label="Порогов"
          value={formatNumber(thresholdsQuery.data?.length ?? 0)}
          caption="Configured route thresholds"
          icon={<TuneRounded />}
        />
        <KpiCard
          label="Office ID"
          value={warehouse.office_from_id}
          caption="Первичный бизнес-идентификатор"
          icon={<WarehouseRounded />}
        />
        <KpiCard
          label="Создан"
          value={formatDateTime(warehouse.created_at)}
          caption="Таймстамп из PostgreSQL"
          icon={<LayersRounded />}
        />
      </Box>

      <Box
        sx={{
          display: "grid",
          gap: 2,
          gridTemplateColumns: { xs: "1fr", xl: "1.2fr 0.8fr" },
        }}
      >
        <GlassPanel>
          <Stack spacing={2}>
            <Typography variant="h5">Маршруты склада</Typography>
            {routesQuery.data?.length ? (
              <TableContainer>
                <Table>
                  <TableHead>
                    <TableRow>
                      <TableCell>Route UUID</TableCell>
                      <TableCell>Business route_id</TableCell>
                      <TableCell>Название</TableCell>
                      <TableCell>Порог</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {routesQuery.data.map((route) => (
                      <TableRow key={route.id}>
                        <TableCell>{route.id}</TableCell>
                        <TableCell>{route.route_id}</TableCell>
                        <TableCell>{route.name || "Без названия"}</TableCell>
                        <TableCell>{thresholdsMap.get(route.id)?.value ?? "—"}</TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </TableContainer>
            ) : (
              <EmptyState
                title="Маршрутов пока нет"
                description="Добавьте первый маршрут, чтобы начать настройку порогов вызова."
                actionLabel="Добавить маршрут"
                onAction={() => setRouteDialogOpen(true)}
              />
            )}
          </Stack>
        </GlassPanel>

        <GlassPanel>
          <Stack spacing={2}>
            <Typography variant="h5">Настройка порога вызова</Typography>
            <Typography color="text.secondary">
              Backend ожидает `route_id` как UUID из таблицы `routes`, не бизнесовый `route_id`.
            </Typography>
            {thresholdError ? <Alert severity="error">{thresholdError}</Alert> : null}
            <TextField
              select
              label="Маршрут"
              value={selectedRouteId}
              onChange={(event) => setSelectedRouteId(event.target.value)}
            >
              {(routesQuery.data ?? []).map((route) => (
                <MenuItem key={route.id} value={route.id}>
                  {route.route_id} {route.name ? `· ${route.name}` : ""}
                </MenuItem>
              ))}
            </TextField>
            <TextField
              label="Порог"
              type="number"
              value={thresholdValue}
              onChange={(event) => setThresholdValue(event.target.value)}
            />
            <Button
              variant="contained"
              onClick={() =>
                updateThresholdMutation.mutate({
                  warehouse_id: warehouseId,
                  route_id: selectedRouteId,
                  value: Number(thresholdValue),
                })
              }
              disabled={!selectedRouteId || !thresholdValue || updateThresholdMutation.isPending}
            >
              {updateThresholdMutation.isPending ? "Сохраняем..." : "Сохранить порог"}
            </Button>
          </Stack>
        </GlassPanel>
      </Box>

      <Dialog open={routeDialogOpen} onClose={() => setRouteDialogOpen(false)} fullWidth maxWidth="sm">
        <DialogTitle>Добавить маршрут</DialogTitle>
        <DialogContent>
          <Stack spacing={2} sx={{ pt: 1 }}>
            {routeError ? <Alert severity="error">{routeError}</Alert> : null}
            <TextField label="Business route_id" value={routeId} onChange={(e) => setRouteId(e.target.value)} />
            <TextField label="Название" value={routeName} onChange={(e) => setRouteName(e.target.value)} />
          </Stack>
        </DialogContent>
        <DialogActions sx={{ px: 3, pb: 3 }}>
          <Button onClick={() => setRouteDialogOpen(false)}>Отмена</Button>
          <Button
            variant="contained"
            disabled={!routeId || createRouteMutation.isPending}
            onClick={() => {
              setRouteError("");
              createRouteMutation.mutate({ route_id: routeId, name: routeName });
            }}
          >
            {createRouteMutation.isPending ? "Создаём..." : "Создать"}
          </Button>
        </DialogActions>
      </Dialog>
    </Stack>
  );
}
