import TrendingUpRounded from "@mui/icons-material/TrendingUpRounded";
import { alpha, Box, Stack, Typography, useTheme } from "@mui/material";
import type { ReactNode } from "react";

import { GlassPanel } from "./GlassPanel";

interface KpiCardProps {
  label: string;
  value: string;
  caption: string;
  icon?: ReactNode;
  accent?: string;
}

export function KpiCard(props: KpiCardProps) {
  const theme = useTheme();
  const accent = props.accent ?? theme.palette.primary.main;

  return (
    <GlassPanel
      sx={{
        minHeight: 168,
        display: "flex",
        flexDirection: "column",
        justifyContent: "space-between",
      }}
    >
      <Stack direction="row" justifyContent="space-between" alignItems="flex-start">
        <Stack spacing={1}>
          <Typography color="text.secondary" variant="body2">
            {props.label}
          </Typography>
          <Typography variant="h4">{props.value}</Typography>
        </Stack>
        <Box
          sx={{
            width: 48,
            height: 48,
            borderRadius: 3,
            display: "grid",
            placeItems: "center",
            color: accent,
            backgroundColor: alpha(accent, 0.18),
          }}
        >
          {props.icon ?? <TrendingUpRounded />}
        </Box>
      </Stack>
      <Typography color="text.secondary">{props.caption}</Typography>
    </GlassPanel>
  );
}
