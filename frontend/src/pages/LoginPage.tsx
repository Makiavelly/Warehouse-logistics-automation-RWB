import KeyRounded from "@mui/icons-material/KeyRounded";
import LockOpenRounded from "@mui/icons-material/LockOpenRounded";
import { Alert, Box, Button, InputAdornment, Stack, TextField, Typography } from "@mui/material";
import { useMutation } from "@tanstack/react-query";
import { useState } from "react";
import { useNavigate } from "react-router-dom";

import { authApi } from "../api/auth";
import { extractApiError } from "../api/client";
import { useSession } from "../app/session";
import { GlassPanel } from "../components/GlassPanel";

export function LoginPage() {
  const navigate = useNavigate();
  const { login } = useSession();
  const [username, setUsername] = useState("admin");
  const [password, setPassword] = useState("admin123");
  const [errorText, setErrorText] = useState("");

  const loginMutation = useMutation({
    mutationFn: authApi.login,
    onSuccess: (response) => {
      login(response.token, response.role);
      navigate("/");
    },
    onError: (error) => {
      setErrorText(extractApiError(error));
    },
  });

  return (
    <Box
      sx={{
        minHeight: "100vh",
        display: "grid",
        placeItems: "center",
        px: 2,
      }}
    >
      <GlassPanel sx={{ width: "100%", maxWidth: 420 }}>
        <Stack
          component="form"
          spacing={2}
          onSubmit={(event) => {
            event.preventDefault();
            setErrorText("");
            loginMutation.mutate({ username, password });
          }}
        >
          <Typography variant="h4">Вход</Typography>
          <Typography color="text.secondary">Авторизация через `POST /api/auth/login`.</Typography>

          <Alert severity="info">
            Demo preset: `admin / admin123`.
          </Alert>

          {errorText ? <Alert severity="error">{errorText}</Alert> : null}

          <TextField
            label="Username"
            value={username}
            onChange={(event) => setUsername(event.target.value)}
            InputProps={{
              startAdornment: (
                <InputAdornment position="start">
                  <KeyRounded />
                </InputAdornment>
              ),
            }}
          />
          <TextField
            label="Password"
            type="password"
            value={password}
            onChange={(event) => setPassword(event.target.value)}
            InputProps={{
              startAdornment: (
                <InputAdornment position="start">
                  <LockOpenRounded />
                </InputAdornment>
              ),
            }}
          />

          <Button type="submit" variant="contained" size="large" disabled={loginMutation.isPending}>
            {loginMutation.isPending ? "Входим..." : "Войти"}
          </Button>
        </Stack>
      </GlassPanel>
    </Box>
  );
}
