import React, { useState, useEffect } from "react";
import { BrowserRouter as Router, Routes, Route, useNavigate, Navigate } from "react-router-dom";
import AuthTabs from "./components/AuthTabs";
import Games from "./pages/Games";
import Bets from "./pages/Bets";
import Profile from "./profile/Profile";

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    const savedAuthStatus = localStorage.getItem("isAuthenticated");
    // Make sure user is logged out on page load
    if (savedAuthStatus === "true") {
      setIsAuthenticated(true);
    }
  }, [navigate]);

  const handleLogin = () => {
    setIsAuthenticated(true);
    localStorage.setItem("isAuthenticated", "true");
    navigate("/profile");
  };

  const handleLogout = () => {
    setIsAuthenticated(false);
    localStorage.removeItem("isAuthenticated");
    navigate("/login");
  };

  return (
    <div className="min-h-screen bg-blue-600 flex items-center justify-center p-4">
      <div className="w-full max-w-md bg-[#B4E4E4] rounded-xl shadow-xl overflow-hidden relative min-h-[600px]">
        {!isAuthenticated ? (
          <AuthTabs onLogin={handleLogin} />
        ) : (
          <Routes>
            <Route path="/login" element={<Profile onLogout={handleLogout} />} />
            <Route path="/games" element={<Games />} />
            <Route path="/bets" element={<Bets />} />
            {/* Redirect to login page if not authenticated */}
            <Route path="*" element={<Navigate to="/login" />} />
          </Routes>
        )}
      </div>
    </div>
  );
}

export default function AppWrapper() {
  return (
    <Router>
      <App />
    </Router>
  );
}
