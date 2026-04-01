import { Paper, type PaperProps } from "@mui/material";

export function GlassPanel(props: PaperProps) {
  return (
    <Paper
      {...props}
      sx={{
        p: { xs: 2.25, md: 3 },
        borderRadius: 6,
        background:
          "linear-gradient(180deg, rgba(255,255,255,0.06) 0%, rgba(255,255,255,0.02) 100%)",
        ...props.sx,
      }}
    />
  );
}
