import { createContext, useContext, useMemo, useState, type PropsWithChildren } from "react";

import {
  clearStoredRole,
  clearStoredToken,
  getStoredApiKey,
  getStoredRole,
  getStoredToken,
  setStoredApiKey,
  setStoredRole,
  setStoredToken,
} from "../utils/storage";

interface SessionContextValue {
  token: string;
  role: string;
  apiKey: string;
  login: (token: string, role: string) => void;
  logout: () => void;
  setApiKey: (apiKey: string) => void;
}

const SessionContext = createContext<SessionContextValue | null>(null);

export function SessionProvider({ children }: PropsWithChildren) {
  const [token, setToken] = useState(getStoredToken);
  const [role, setRole] = useState(getStoredRole);
  const [apiKey, setApiKeyState] = useState(getStoredApiKey);

  const value = useMemo<SessionContextValue>(
    () => ({
      token,
      role,
      apiKey,
      login(nextToken, nextRole) {
        setToken(nextToken);
        setRole(nextRole);
        setStoredToken(nextToken);
        setStoredRole(nextRole);
      },
      logout() {
        setToken("");
        setRole("");
        clearStoredToken();
        clearStoredRole();
      },
      setApiKey(nextApiKey) {
        setApiKeyState(nextApiKey);
        setStoredApiKey(nextApiKey);
      },
    }),
    [apiKey, role, token],
  );

  return <SessionContext.Provider value={value}>{children}</SessionContext.Provider>;
}

export function useSession() {
  const context = useContext(SessionContext);

  if (!context) {
    throw new Error("useSession must be used within SessionProvider");
  }

  return context;
}
