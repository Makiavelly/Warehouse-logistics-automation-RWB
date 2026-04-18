import { BrowserRouter, Navigate, Outlet, Route, Routes } from "react-router-dom";

import { AppShell } from "./components/AppShell";
import { useSession } from "./app/session";
import { AnalyticsPage } from "./pages/AnalyticsPage";
import { DashboardPage } from "./pages/DashboardPage";
import { ForecastsPage } from "./pages/ForecastsPage";
import { HorizonsPage } from "./pages/HorizonsPage";
import { LoginPage } from "./pages/LoginPage";
import { NotificationsPage } from "./pages/NotificationsPage";
import { SettingsPage } from "./pages/SettingsPage";
import { WarehousesPage } from "./pages/WarehousesPage";
import { WarehouseDetailsPage } from "./pages/WarehouseDetailsPage";

function RequireAuth() {
  const { token } = useSession();

  if (!token) {
    return <Navigate to="/login" replace />;
  }

  return (
    <AppShell>
      <Outlet />
    </AppShell>
  );
}

function LoginGate() {
  const { token } = useSession();

  if (token) {
    return <Navigate to="/" replace />;
  }

  return <LoginPage />;
}

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<LoginGate />} />
        <Route element={<RequireAuth />}>
          <Route path="/" element={<DashboardPage />} />
          <Route path="/warehouses" element={<WarehousesPage />} />
          <Route path="/warehouses/:warehouseId" element={<WarehouseDetailsPage />} />
          <Route path="/forecasts" element={<ForecastsPage />} />
          <Route path="/notifications" element={<NotificationsPage />} />
          <Route path="/horizons" element={<HorizonsPage />} />
          <Route path="/analytics" element={<AnalyticsPage />} />
          <Route path="/settings" element={<SettingsPage />} />
        </Route>
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
}
