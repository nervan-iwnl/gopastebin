// frontend/src/App.jsx
import { Routes, Route, Navigate } from "react-router-dom";
import Login from "./pages/Login.jsx";
import Register from "./pages/Register.jsx";
import MyPastes from "./pages/MyPastes.jsx";
import PasteView from "./pages/PasteView.jsx";
import AdminStorage from "./pages/AdminStorage.jsx";
import NotFound from "./pages/NotFound.jsx";
import Nav from "./components/Nav.jsx";
import { useAuth } from "./hooks/useAuth.jsx";

function PrivateRoute({ children }) {
  const { user, loading } = useAuth();
  if (loading) return <div>loading...</div>;
  if (!user) return <Navigate to="/login" replace />;
  return children;
}

function AdminRoute({ children }) {
  const { user, loading } = useAuth();
  if (loading) return <div>loading...</div>;
  if (!user || !user.is_admin) return <Navigate to="/" replace />;
  return children;
}

function App() {
  return (
    <div className="app">
      <Nav />
      <div className="container" style={{ maxWidth: 900, margin: "20px auto" }}>
        <Routes>
          <Route
            path="/"
            element={
              <PrivateRoute>
                <MyPastes />
              </PrivateRoute>
            }
          />
          <Route
            path="/p/:slug"
            element={
              // allow public viewing of a paste by link (no auth required)
              <PasteView />
            }
          />
          <Route
            path="/admin/storage"
            element={
              <AdminRoute>
                <AdminStorage />
              </AdminRoute>
            }
          />

          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />

          <Route path="*" element={<NotFound />} />
        </Routes>
      </div>
    </div>
  );
}

export default App;
