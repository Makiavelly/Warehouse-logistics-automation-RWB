import { Chip, Stack, Typography } from "@mui/material";

interface PageHeaderProps {
  eyebrow?: string;
  title: string;
  description: string;
  chipLabel?: string;
  chipColor?: "default" | "primary" | "secondary" | "success" | "warning" | "error";
}

export function PageHeader(props: PageHeaderProps) {
  const { eyebrow, title, description, chipLabel, chipColor = "secondary" } = props;

  return (
    <Stack spacing={1.25}>
      <Stack
        direction={{ xs: "column", sm: "row" }}
        spacing={1.25}
        alignItems={{ xs: "flex-start", sm: "center" }}
      >
        {eyebrow ? (
          <Typography
            variant="overline"
            sx={{ letterSpacing: "0.18em", color: "text.secondary", fontWeight: 700 }}
          >
            {eyebrow}
          </Typography>
        ) : null}
        {chipLabel ? <Chip size="small" color={chipColor} label={chipLabel} /> : null}
      </Stack>
      <Typography variant="h3">{title}</Typography>
      <Typography color="text.secondary" sx={{ maxWidth: 820 }}>
        {description}
      </Typography>
    </Stack>
  );
}
