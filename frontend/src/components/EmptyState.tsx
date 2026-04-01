import { Button, Stack, Typography } from "@mui/material";
import InboxOutlined from "@mui/icons-material/InboxOutlined";

import { GlassPanel } from "./GlassPanel";

interface EmptyStateProps {
  title: string;
  description: string;
  actionLabel?: string;
  onAction?: () => void;
}

export function EmptyState(props: EmptyStateProps) {
  const { title, description, actionLabel, onAction } = props;

  return (
    <GlassPanel>
      <Stack spacing={2} alignItems="flex-start">
        <InboxOutlined color="primary" />
        <Typography variant="h5">{title}</Typography>
        <Typography color="text.secondary">{description}</Typography>
        {actionLabel && onAction ? (
          <Button variant="contained" onClick={onAction}>
            {actionLabel}
          </Button>
        ) : null}
      </Stack>
    </GlassPanel>
  );
}
