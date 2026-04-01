import { Alert, Box, Button, Stack, Typography } from "@mui/material";
import { Component, type ErrorInfo, type PropsWithChildren, type ReactNode } from "react";

import { GlassPanel } from "./GlassPanel";

interface State {
  hasError: boolean;
}

export class AppErrorBoundary extends Component<PropsWithChildren, State> {
  override state: State = {
    hasError: false,
  };

  override componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error("AppErrorBoundary", error, errorInfo);
  }

  static getDerivedStateFromError() {
    return { hasError: true };
  }

  override render(): ReactNode {
    if (this.state.hasError) {
      return (
        <Box sx={{ minHeight: "100vh", display: "grid", placeItems: "center", p: 3 }}>
          <GlassPanel sx={{ maxWidth: 520, width: "100%" }}>
            <Stack spacing={2}>
              <Typography variant="h4">Интерфейс поймал критическую ошибку</Typography>
              <Alert severity="error">
                Проверьте консоль браузера и перезагрузите страницу. Если ошибка повторится,
                проблема в одном из фронтенд-модулей.
              </Alert>
              <Button variant="contained" onClick={() => window.location.reload()}>
                Перезагрузить приложение
              </Button>
            </Stack>
          </GlassPanel>
        </Box>
      );
    }

    return this.props.children;
  }
}
