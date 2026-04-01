import { Chip } from "@mui/material";

const horizonMeta = {
  recommended: { label: "Recommended", color: "success" },
  acceptable: { label: "Acceptable", color: "warning" },
  use_with_caution: { label: "Caution", color: "warning" },
  unreliable: { label: "Unreliable", color: "error" },
} as const;

export function HorizonTag({ level }: { level: keyof typeof horizonMeta }) {
  const meta = horizonMeta[level];
  return <Chip size="small" color={meta.color} label={meta.label} />;
}
