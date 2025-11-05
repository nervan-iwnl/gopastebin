import React, { createContext, useContext, useEffect, useState } from "react";
import { api } from "../api";

const AuthCtx = createContext(null);

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  // грузим /auth/me при старте
  useEffect(() => {
    (async () => {
      try {
        const data = await api.me();
        setUser(data.user);
      } catch {
        // если /auth/me не доступен (например CORS/cookie проблемы),
        // попытаться восстановить user из localStorage как graceful fallback
        try {
          const saved = localStorage.getItem("gopaste_user");
          if (saved) setUser(JSON.parse(saved));
          else setUser(null);
        } catch {
          setUser(null);
        }
      } finally {
        setLoading(false);
      }
    })();
  }, []);

  async function login(login, password) {
    // вызвать login на бэке, затем попытаться получить /me
    const res = await api.login(login, password);
    // if backend returned tokens in body, store them for Authorization header
    try {
      if (res?.access) {
        localStorage.setItem("gopaste_token", res.access);
      }
      if (res?.refresh) {
        localStorage.setItem("gopaste_refresh", res.refresh);
      }
    } catch {}

    // если login вернул user — используем его, иначе дергаем /me
    if (res?.user) {
      setUser(res.user);
      try {
        localStorage.setItem("gopaste_user", JSON.stringify(res.user));
      } catch {}
      return;
    }

    const data = await api.me();
    setUser(data.user);
    try {
      localStorage.setItem("gopaste_user", JSON.stringify(data.user));
    } catch {}
  }

  async function register(email, username, password) {
    await api.register(email, username, password);
    // после регистрации не факт что верифицировался — оставим как есть
  }

  function logout() {
    // бэк не даёт logout — просто чистим
    setUser(null);
    try {
      localStorage.removeItem("gopaste_user");
      localStorage.removeItem("gopaste_token");
      localStorage.removeItem("gopaste_refresh");
    } catch {}
  }

  return (
    <AuthCtx.Provider value={{ user, loading, login, register, logout, setUser }}>
      {children}
    </AuthCtx.Provider>
  );
}

export function useAuth() {
  return useContext(AuthCtx);
}
