import "@fontsource/manrope/400.css";
import "@fontsource/manrope/500.css";
import "@fontsource/manrope/700.css";
import "@fontsource/space-grotesk/500.css";
import "@fontsource/space-grotesk/700.css";
import React from "react";
import ReactDOM from "react-dom/client";
import { CssBaseline, ThemeProvider } from "@mui/material";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

import App from "./App";
import { AppErrorBoundary } from "./components/AppErrorBoundary";
import { SessionProvider } from "./app/session";
import { appTheme } from "./theme";
import "./index.css";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
      staleTime: 30_000,
    },
  },
});

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <AppErrorBoundary>
      <QueryClientProvider client={queryClient}>
        <SessionProvider>
          <ThemeProvider theme={appTheme}>
            <CssBaseline />
            <App />
          </ThemeProvider>
        </SessionProvider>
      </QueryClientProvider>
    </AppErrorBoundary>
  </React.StrictMode>,
);
