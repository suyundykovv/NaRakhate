import React, { useState } from "react"
import { BrowserRouter as Router, Routes, Route } from "react-router-dom"
import AuthTabs from "./components/AuthTabs"

import Games from "./pages/Games"
import Bets from "./pages/Bets"
import Profile from "./pages/Profile"

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false)

  const handleLogin = () => {
    setIsAuthenticated(true)
  }

  return (
    <Router>
      <div className="min-h-screen bg-blue-600 flex items-center justify-center p-4">
        <div className="w-full max-w-md bg-[#B4E4E4] rounded-xl shadow-xl overflow-hidden relative min-h-[600px]">
          {!isAuthenticated ? (
            <AuthTabs onLogin={handleLogin} />
          ) : (
            <>
              <Routes>

                <Route path="/games" element={<Games />} />
                <Route path="/bets" element={<Bets />} />
                <Route path="/profile" element={<Profile />} />
                <Route path="/" element={<Profile />} />
              </Routes>

            </>
          )}
        </div>
      </div>
    </Router>
  )
}

export default App
