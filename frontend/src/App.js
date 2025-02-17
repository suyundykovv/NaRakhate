import React, { useState, useEffect } from "react";
import { BrowserRouter as Router, Routes, Route, useNavigate } from "react-router-dom";
import AuthTabs from "./components/AuthTabs";
import Games from "./pages/Games";
import Bets from "./pages/Bets";
import Profile from "./profile/Profile";

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const navigate = useNavigate();

  // Check localStorage for authentication status on initial load
  useEffect(() => {
    const savedAuthStatus = localStorage.getItem("isAuthenticated");
    if (savedAuthStatus === "true") {
      setIsAuthenticated(true);
    }
  }, [navigate]);

  const handleLogin = () => {
    setIsAuthenticated(true);
    localStorage.setItem("isAuthenticated", "true"); // Save authentication status to localStorage
    navigate("/profile"); // Redirect to profile after login
  };

  const handleLogout = () => {
    setIsAuthenticated(false);
    localStorage.removeItem("isAuthenticated"); // Clear authentication status from localStorage
    navigate("/login"); // Redirect to login after logout
  };

  return (
    <div className="min-h-screen bg-blue-600 flex items-center justify-center p-4">
      <div className="w-full max-w-md bg-[#B4E4E4] rounded-xl shadow-xl overflow-hidden relative min-h-[600px]">
        {!isAuthenticated ? (
          <AuthTabs onLogin={handleLogin} />
        ) : (
          <Routes>
            <Route path="/profile" element={<Profile onLogout={handleLogout} />} />
            <Route path="/games" element={<Games />} />
            <Route path="/bets" element={<Bets />} />
          </Routes>
        )}
      </div>
    </div>
  );
}

// Wrap App with Router
export default function AppWrapper() {
  return (
    <Router>
      <App />
    </Router>
  );
}