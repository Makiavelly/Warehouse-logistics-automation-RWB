const TOKEN_KEY = "wb-session-token";
const ROLE_KEY = "wb-session-role";
const API_KEY = "wb-ingest-api-key";

export function getStoredToken() {
  return localStorage.getItem(TOKEN_KEY) ?? "";
}

export function setStoredToken(token: string) {
  localStorage.setItem(TOKEN_KEY, token);
}

export function clearStoredToken() {
  localStorage.removeItem(TOKEN_KEY);
}

export function getStoredRole() {
  return localStorage.getItem(ROLE_KEY) ?? "";
}

export function setStoredRole(role: string) {
  localStorage.setItem(ROLE_KEY, role);
}

export function clearStoredRole() {
  localStorage.removeItem(ROLE_KEY);
}

export function getStoredApiKey() {
  return localStorage.getItem(API_KEY) ?? "internal-api-key";
}

export function setStoredApiKey(apiKey: string) {
  localStorage.setItem(API_KEY, apiKey);
}
