import AddRounded from "@mui/icons-material/AddRounded";
import {
  Alert,
  Box,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  LinearProgress,
  Stack,
  TextField,
  Typography,
} from "@mui/material";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useMemo, useState } from "react";
import { Link as RouterLink } from "react-router-dom";

import { extractApiError } from "../api/client";
import { thresholdsApi } from "../api/thresholds";
import { warehousesApi } from "../api/warehouses";
import { EmptyState } from "../components/EmptyState";
import { GlassPanel } from "../components/GlassPanel";
import { LoadingBlock } from "../components/LoadingBlock";
import { formatDate } from "../utils/format";
import type { Warehouse } from "../types/models";

function WarehouseCard({ warehouse, thresholdCount }: { warehouse: Warehouse; thresholdCount: number }) {
  const routesQuery = useQuery({
    queryKey: ["warehouse-card-routes", warehouse.id],
    queryFn: () => warehousesApi.getRoutes(warehouse.id),
  });

  const routesCount = routesQuery.data?.length ?? 0;
  const coverage = routesCount ? Math.min(100, Math.round((thresholdCount / routesCount) * 100)) : 0;

  return (
    <GlassPanel sx={{ minHeight: 280 }}>
      <Stack spacing={2.25} sx={{ height: "100%" }}>
        <Stack spacing={0.75}>
          <Typography variant="h5">{warehouse.name}</Typography>
          <Typography color="text.secondary">{warehouse.office_from_id}</Typography>
        </Stack>

        <Typography color="text.secondary">
          {warehouse.address || "Адрес пока не указан. Можно завести позже."}
        </Typography>

        <Stack direction="row" spacing={2}>
          <Stack>
            <Typography variant="h6">{routesCount}</Typography>
            <Typography color="text.secondary" variant="body2">
              маршрутов
            </Typography>
          </Stack>
          <Stack>
            <Typography variant="h6">{thresholdCount}</Typography>
            <Typography color="text.secondary" variant="body2">
              порогов
            </Typography>
          </Stack>
        </Stack>

        <Box>
          <Stack direction="row" justifyContent="space-between" sx={{ mb: 0.75 }}>
            <Typography variant="body2">Готовность конфигурации</Typography>
            <Typography variant="body2" color="text.secondary">
              {coverage}%
            </Typography>
          </Stack>
          <LinearProgress variant="determinate" value={coverage} sx={{ height: 10, borderRadius: 999 }} />
        </Box>

        <Stack direction="row" justifyContent="space-between" alignItems="center" sx={{ mt: "auto" }}>
          <Typography color="text.secondary" variant="body2">
            Создан {formatDate(warehouse.created_at)}
          </Typography>
          <Button component={RouterLink} to={`/warehouses/${warehouse.id}`} variant="contained">
            Подробнее
          </Button>
        </Stack>
      </Stack>
    </GlassPanel>
  );
}

export function WarehousesPage() {
  const queryClient = useQueryClient();
  const [search, setSearch] = useState("");
  const [open, setOpen] = useState(false);
  const [name, setName] = useState("");
  const [officeFromId, setOfficeFromId] = useState("");
  const [address, setAddress] = useState("");
  const [formError, setFormError] = useState("");

  const warehousesQuery = useQuery({
    queryKey: ["warehouses"],
    queryFn: warehousesApi.getAll,
  });

  const thresholdsQuery = useQuery({
    queryKey: ["thresholds"],
    queryFn: () => thresholdsApi.getAll(),
  });

  const createWarehouseMutation = useMutation({
    mutationFn: warehousesApi.create,
    onSuccess: async () => {
      setOpen(false);
      setName("");
      setOfficeFromId("");
      setAddress("");
      await queryClient.invalidateQueries({ queryKey: ["warehouses"] });
    },
    onError: (error) => {
      setFormError(extractApiError(error));
    },
  });

  const filteredWarehouses = useMemo(() => {
    const needle = search.trim().toLowerCase();

    if (!needle) {
      return warehousesQuery.data ?? [];
    }

    return (warehousesQuery.data ?? []).filter((warehouse) =>
      [warehouse.name, warehouse.office_from_id, warehouse.address]
        .join(" ")
        .toLowerCase()
        .includes(needle),
    );
  }, [search, warehousesQuery.data]);

  return (
    <Stack spacing={3}>
      <GlassPanel>
        <Stack
          direction={{ xs: "column", md: "row" }}
          spacing={2}
          justifyContent="space-between"
          alignItems={{ xs: "stretch", md: "center" }}
        >
          <TextField
            label="Поиск по складам"
            value={search}
            onChange={(event) => setSearch(event.target.value)}
            placeholder="Название, office_from_id или адрес"
          />
          <Button startIcon={<AddRounded />} variant="contained" onClick={() => setOpen(true)}>
            Добавить склад
          </Button>
        </Stack>
      </GlassPanel>

      {warehousesQuery.isLoading ? (
        <LoadingBlock />
      ) : filteredWarehouses.length ? (
        <Box
          sx={{
            display: "grid",
            gap: 2,
            gridTemplateColumns: { xs: "1fr", lg: "repeat(2, 1fr)", xl: "repeat(3, 1fr)" },
          }}
        >
          {filteredWarehouses.map((warehouse) => (
            <WarehouseCard
              key={warehouse.id}
              warehouse={warehouse}
              thresholdCount={
                thresholdsQuery.data?.filter((item) => item.warehouse_id === warehouse.id).length ?? 0
              }
            />
          ))}
        </Box>
      ) : (
        <EmptyState
          title="Склады не найдены"
          description="Создайте первый склад или измените строку поиска."
          actionLabel="Создать склад"
          onAction={() => setOpen(true)}
        />
      )}

      <Dialog open={open} onClose={() => setOpen(false)} fullWidth maxWidth="sm">
        <DialogTitle>Новый склад</DialogTitle>
        <DialogContent>
          <Stack spacing={2} sx={{ pt: 1 }}>
            {formError ? <Alert severity="error">{formError}</Alert> : null}
            <TextField label="Название склада" value={name} onChange={(e) => setName(e.target.value)} />
            <TextField
              label="Office From ID"
              value={officeFromId}
              onChange={(e) => setOfficeFromId(e.target.value)}
            />
            <TextField label="Адрес" value={address} onChange={(e) => setAddress(e.target.value)} />
          </Stack>
        </DialogContent>
        <DialogActions sx={{ px: 3, pb: 3 }}>
          <Button variant="text" onClick={() => setOpen(false)}>
            Отмена
          </Button>
          <Button
            variant="contained"
            onClick={() => {
              setFormError("");
              createWarehouseMutation.mutate({
                name,
                office_from_id: officeFromId,
                address,
              });
            }}
            disabled={!name || !officeFromId || createWarehouseMutation.isPending}
          >
            {createWarehouseMutation.isPending ? "Создаём..." : "Создать"}
          </Button>
        </DialogActions>
      </Dialog>
    </Stack>
  );
}
