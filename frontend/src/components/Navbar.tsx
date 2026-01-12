import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

function Navbar() {
  const { user, logout } = useAuth();

  return (
    <nav className="navbar">
      <Link to="/" className="navbar-brand">
        � PollHub
      </Link>
      <div className="navbar-nav">
        {user ? (
          <>
            <span className="navbar-user">Welcome, {user.username}</span>
            <button onClick={logout} className="btn btn-outline">
              Logout
            </button>
          </>
        ) : (
          <>
            <Link to="/login" className="btn btn-outline">
              Login
            </Link>
            <Link to="/signup" className="btn btn-primary">
              Sign Up
            </Link>
          </>
        )}
      </div>
    </nav>
  );
}

export default Navbar;
