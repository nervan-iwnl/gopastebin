// frontend/src/api.js
const API_BASE =
  import.meta.env.VITE_API_BASE || "http://localhost:8080/api/v1";

async function request(path, { method = "GET", body, headers } = {}) {
  // attach Authorization header from saved token if present
  const savedToken = (() => {
    try {
      return localStorage.getItem("gopaste_token");
    } catch {
      return null;
    }
  })();

  const reqHeaders = {
    "Content-Type": "application/json",
    ...(headers || {}),
  };
  if (savedToken) reqHeaders["Authorization"] = `Bearer ${savedToken}`;

  const res = await fetch(`${API_BASE}${path}`, {
    method,
    headers: reqHeaders,
    body: body ? JSON.stringify(body) : undefined,
    credentials: "include", // still include cookies when available
  });

  const text = await res.text();
  let data;
  try {
    data = text ? JSON.parse(text) : {};
  } catch {
    data = { raw: text };
  }

  if (!res.ok) {
    throw data;
  }
  return data;
}

// для /raw нам нужен просто текст
async function raw(path) {
  const res = await fetch(`${API_BASE}${path}`, {
    credentials: "include",
  });
  const text = await res.text();
  if (!res.ok) {
    throw { error: { message: text || "failed" } };
  }
  return text;
}

export const api = {
  me: () => request("/auth/me"),
  login: (login, password) =>
    request("/auth/login", { method: "POST", body: { login, password } }),
  register: (email, username, password) =>
    request("/auth/register", {
      method: "POST",
      body: { email, username, password },
    }),

  myPastes: () => request("/me/pastes"),
  createPaste: (payload) =>
    request("/pastes", { method: "POST", body: payload }),
  createAnon: (content) =>
    request("/pastes/anon", { method: "POST", body: { title: "anon", content, extension: "txt", folder: "", is_public: true } }),
  getPaste: (slug) => request(`/pastes/${slug}`),
  getPasteRaw: (slug) => raw(`/pastes/${slug}/raw`),
  getRecent: () => request("/pastes/recent"),

  getStorage: () => request("/settings/storage"),
  setStorage: (storage) =>
    request("/settings/storage", { method: "POST", body: { storage } }),
};
