import { Skeleton, Stack } from "@mui/material";

import { GlassPanel } from "./GlassPanel";

export function LoadingBlock() {
  return (
    <GlassPanel>
      <Stack spacing={1.5}>
        <Skeleton variant="text" width="28%" height={36} />
        <Skeleton variant="rounded" width="100%" height={68} />
        <Skeleton variant="rounded" width="100%" height={68} />
        <Skeleton variant="rounded" width="84%" height={68} />
      </Stack>
    </GlassPanel>
  );
}
