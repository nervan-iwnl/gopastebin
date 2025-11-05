// frontend/src/components/Nav.jsx
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../hooks/useAuth.jsx";

export default function Nav() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  function handleLogout() {
    logout();
    navigate("/login");
  }

  return (
    <nav className="navbar navbar-dark bg-dark">
      <div className="container">
        <Link className="navbar-brand" to="/">
          gopastebin
        </Link>

        <div className="d-flex align-items-center">
          <Link className="nav-link text-white me-3" to="/">
            My pastes
          </Link>

          {user?.is_admin && (
            <Link className="nav-link text-white me-3" to="/admin/storage">
              Admin
            </Link>
          )}

          {user ? (
            <div className="d-flex align-items-center">
              <span className="text-white me-2">{user.username}</span>
              <button className="btn btn-sm btn-outline-light" onClick={handleLogout}>
                Logout
              </button>
            </div>
          ) : (
            <div>
              <Link className="btn btn-sm btn-outline-light me-2" to="/login">
                Login
              </Link>
              <Link className="btn btn-sm btn-light" to="/register">
                Register
              </Link>
            </div>
          )}
        </div>
      </div>
    </nav>
  );
}
