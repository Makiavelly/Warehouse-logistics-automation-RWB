import { Chip } from "@mui/material";

const statusMeta = {
  pending: { label: "Pending", color: "warning" },
  sent: { label: "Sent", color: "info" },
  confirmed: { label: "Confirmed", color: "success" },
  completed: { label: "Completed", color: "success" },
  cancelled: { label: "Cancelled", color: "error" },
  accepted: { label: "Accepted", color: "success" },
  missed: { label: "Missed", color: "error" },
} as const;

export function StatusTag({ status }: { status: keyof typeof statusMeta }) {
  const meta = statusMeta[status];
  return <Chip size="small" color={meta.color} label={meta.label} />;
}
