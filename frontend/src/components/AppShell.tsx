import MenuRounded from "@mui/icons-material/MenuRounded";
import LogoutRounded from "@mui/icons-material/LogoutRounded";
import AnalyticsRounded from "@mui/icons-material/AnalyticsRounded";
import DashboardRounded from "@mui/icons-material/DashboardRounded";
import InsightsRounded from "@mui/icons-material/InsightsRounded";
import Inventory2Rounded from "@mui/icons-material/Inventory2Rounded";
import NotificationsActiveRounded from "@mui/icons-material/NotificationsActiveRounded";
import SettingsRounded from "@mui/icons-material/SettingsRounded";
import TrackChangesRounded from "@mui/icons-material/TrackChangesRounded";
import LocalShippingRounded from "@mui/icons-material/LocalShippingRounded";
import {
  Box,
  Drawer,
  IconButton,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Stack,
  Toolbar,
  Typography,
  useMediaQuery,
  useTheme,
} from "@mui/material";
import { useState, type PropsWithChildren } from "react";
import { Link as RouterLink, useLocation, useNavigate } from "react-router-dom";

import { useSession } from "../app/session";

const drawerWidth = 298;

const navigation = [
  { label: "Дашборд", to: "/", icon: <DashboardRounded /> },
  { label: "Склады", to: "/warehouses", icon: <Inventory2Rounded /> },
  { label: "Прогнозы", to: "/forecasts", icon: <InsightsRounded /> },
  { label: "Уведомления", to: "/notifications", icon: <NotificationsActiveRounded /> },
  { label: "Горизонты", to: "/horizons", icon: <TrackChangesRounded /> },
  { label: "Аналитика", to: "/analytics", icon: <AnalyticsRounded /> },
  { label: "Настройки", to: "/settings", icon: <SettingsRounded /> },
];

export function AppShell({ children }: PropsWithChildren) {
  const theme = useTheme();
  const location = useLocation();
  const navigate = useNavigate();
  const { logout } = useSession();
  const [mobileOpen, setMobileOpen] = useState(false);
  const isDesktop = useMediaQuery(theme.breakpoints.up("lg"));

  const drawerContent = (
    <Stack sx={{ height: "100%", p: 2.5 }}>
      <Stack spacing={1.5} sx={{ px: 1, py: 2 }}>
        <Stack direction="row" spacing={1.5} alignItems="center">
          <Box
            sx={{
              width: 44,
              height: 44,
              borderRadius: 3.5,
              display: "grid",
              placeItems: "center",
              background: "linear-gradient(135deg, #CB11AB 0%, #481173 100%)",
              boxShadow: "0 12px 32px rgba(203, 17, 171, 0.24)",
            }}
          >
            <LocalShippingRounded />
          </Box>
          <Stack>
            <Typography variant="h6" sx={{ lineHeight: 1 }}>
              WildFlow Console
            </Typography>
            <Typography color="text.secondary" variant="body2">
              Warehouse logistics automation
            </Typography>
          </Stack>
        </Stack>
      </Stack>

      <List sx={{ py: 1 }}>
        {navigation.map((item) => {
          const selected =
            item.to === "/" ? location.pathname === item.to : location.pathname.startsWith(item.to);

          return (
            <ListItem key={item.to} disablePadding sx={{ mb: 0.75 }}>
              <ListItemButton
                component={RouterLink}
                to={item.to}
                selected={selected}
                onClick={() => setMobileOpen(false)}
                sx={{
                  minHeight: 52,
                  borderRadius: 4,
                  "&.Mui-selected": {
                    background:
                      "linear-gradient(135deg, rgba(203,17,171,0.24), rgba(72,17,115,0.16))",
                  },
                }}
              >
                <ListItemIcon sx={{ minWidth: 40, color: selected ? "primary.main" : "inherit" }}>
                  {item.icon}
                </ListItemIcon>
                <ListItemText primary={item.label} />
              </ListItemButton>
            </ListItem>
          );
        })}
      </List>

      <Box sx={{ mt: "auto", p: 1 }}>
        <Stack spacing={1}>
          <ListItemButton
            onClick={() => {
              logout();
              navigate("/login");
            }}
            sx={{ px: 1, borderRadius: 3 }}
          >
            <ListItemIcon sx={{ minWidth: 38 }}>
              <LogoutRounded />
            </ListItemIcon>
            <ListItemText primary="Выйти" />
          </ListItemButton>
        </Stack>
      </Box>
    </Stack>
  );

  return (
    <Box sx={{ display: "flex", minHeight: "100vh" }}>
      <Box component="nav" sx={{ width: { lg: drawerWidth }, flexShrink: { lg: 0 } }}>
        {!isDesktop ? (
          <Toolbar sx={{ position: "fixed", top: 12, left: 12, zIndex: 1300, minHeight: "auto", p: 0 }}>
            <IconButton
              color="inherit"
              onClick={() => setMobileOpen(true)}
              sx={{
                backdropFilter: "blur(16px)",
                backgroundColor: "rgba(20, 8, 29, 0.72)",
                border: "1px solid rgba(255,255,255,0.08)",
              }}
            >
              <MenuRounded />
            </IconButton>
          </Toolbar>
        ) : null}
        <Drawer
          variant={isDesktop ? "permanent" : "temporary"}
          open={isDesktop ? true : mobileOpen}
          onClose={() => setMobileOpen(false)}
          ModalProps={{ keepMounted: true }}
          sx={{
            "& .MuiDrawer-paper": {
              width: drawerWidth,
              boxSizing: "border-box",
            },
          }}
        >
          {drawerContent}
        </Drawer>
      </Box>

      <Box
        component="main"
        sx={{
          flexGrow: 1,
          px: { xs: 2, md: 3.5 },
          pt: { xs: 3, md: 4 },
          pb: 4,
          width: { lg: `calc(100% - ${drawerWidth}px)` },
        }}
      >
        {children}
      </Box>
    </Box>
  );
}
