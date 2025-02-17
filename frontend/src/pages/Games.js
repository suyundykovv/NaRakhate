import React, { useState, useEffect } from "react";
import {
  ArrowLeft,
  ShoppingCart,
  Layout,
  ChevronRight,
  Gamepad2,
  Receipt,
  UserCircle,
} from "lucide-react";
import { Link } from "react-router-dom";
import "./style_games.css";

function Games() {
  const [activeTab, setActiveTab] = useState("games");
  const [selectedLeague, setSelectedLeague] = useState(null);
  const [events, setEvents] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [selectedBet, setSelectedBet] = useState(null); // Selected bet (match and bet type)
  const [betAmount, setBetAmount] = useState(""); // Bet amount
  const [showBetModal, setShowBetModal] = useState(false); // Show bet modal
  const [userId, setUserId] = useState(null); // Store the user ID

  const fetchUserId = async () => {
    const token = localStorage.getItem("token");
    if (!token) {
      setError("No token found. Please log in again.");
      return;
    }

    try {
      const response = await fetch("http://127.0.0.1:8080/getUser", {
        method: "GET",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) throw new Error("Failed to fetch user data");

      const userData = await response.json();
      setUserId(userData.id); // Set the user ID
    } catch (err) {
      setError("Error fetching user data: " + err.message);
    }
  };

  // Fetch events from the backend
  const fetchEvents = async () => {
    try {
      const response = await fetch("http://127.0.0.1:8080/getEvents");
      if (!response.ok) throw new Error("Failed to load events");
      const data = await response.json();
      return data;
    } catch (error) {
      setError(error.message);
      return [];
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    const loadData = async () => {
      await fetchUserId(); // Fetch user ID first
      const eventsData = await fetchEvents(); // Then fetch events
      setEvents(eventsData);
    };
    loadData();
  }, []);

  // Handle bet button click
  const handleBetClick = (event, oddSelection, oddValue) => {
    setSelectedBet({ event, oddSelection, oddValue });
    setShowBetModal(true);
  };

  const handlePlaceBet = async () => {
    if (!selectedBet || !betAmount || isNaN(betAmount) || betAmount <= 0) {
      alert("Please enter a valid bet amount");
      return;
    }
  
    if (!userId) {
      alert("User ID is missing. Please log in again.");
      return;
    }
  
    const betData = {
      user_id: userId, // Include the user ID
      event_id: selectedBet.event.id,
      odd_selection: selectedBet.oddSelection,
      odd_value: selectedBet.oddValue,
      amount: parseFloat(betAmount),
    };
  
    try {
      const response = await fetch("http://127.0.0.1:8080/createBet", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(betData),
      });
  
      if (!response.ok) {
        const errorText = await response.text(); // Get the error message from the response body
        throw new Error(`Bet placement failed: ${errorText}`);
      }
  
      const data = await response.json();
      alert("Bet placed successfully!");
      setShowBetModal(false);
      setSelectedBet(null);
      setBetAmount("");
    } catch (error) {
      console.error("Error placing bet:", error); // Log detailed error
      alert("Failed to place bet: " + error.message); // Show error message in the UI
    }
  };
  

  return (
    <>
      <div className="container">
        <div className="header">
          <h1>Football Matches</h1>
        </div>

        <div className="content">
          {loading && <div className="loading">Loading matches...</div>}
          {error && <div className="error">Error: {error}</div>}
          {!loading && !error && events.length === 0 && (
            <div className="empty">No matches available</div>
          )}

          <div className="match-list">
            {events.map((event) => (
              <div key={event.id} className="match-item">
                <div className="match-header">
                  <span>
                    {new Date(event.start_time).toLocaleDateString("ru-RU", {
                      day: "numeric",
                      month: "long",
                      hour: "2-digit",
                      minute: "2-digit",
                    })}
                  </span>
                  <div className="match-actions">
                    <ShoppingCart size={18} />
                    <Layout size={18} />
                  </div>
                </div>

                <div className="teams">
                  <div className="team team-left">
                    <span className="team-name">
                      {event.name.split(" vs ")[0]}
                    </span>
                  </div>
                  <div className="vs">VS</div>
                  <div className="team team-right">
                    <span className="team-name">
                      {event.name.split(" vs ")[1]}
                    </span>
                  </div>
                </div>

                <div className="betting-buttons">
                  <button
                    className="bet-button"
                    onClick={() =>
                      handleBetClick(event, "home", event.odds.home_win)
                    }
                  >
                    P1 ({event.odds.home_win.toFixed(2)})
                  </button>
                  <button
                    className="bet-button"
                    onClick={() =>
                      handleBetClick(event, "draw", event.odds.draw)
                    }
                  >
                    N ({event.odds.draw.toFixed(2)})
                  </button>
                  <button
                    className="bet-button"
                    onClick={() =>
                      handleBetClick(event, "away", event.odds.away_win)
                    }
                  >
                    P2 ({event.odds.away_win.toFixed(2)})
                  </button>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Bet Modal */}
        {showBetModal && (
          <div className="bet-modal">
            <div className="bet-modal-content">
              <h3>Place Bet</h3>
              <p>Match: {selectedBet.event.name}</p>
              <p>
                Bet:{" "}
                {selectedBet.oddSelection === "home"
                  ? "P1"
                  : selectedBet.oddSelection === "away"
                  ? "P2"
                  : "N"}
              </p>
              <p>Odds: {selectedBet.oddValue.toFixed(2)}</p>
              <input
                type="number"
                placeholder="Bet amount"
                value={betAmount}
                onChange={(e) => setBetAmount(e.target.value)}
              />
              <button onClick={handlePlaceBet}>Place Bet</button>
              <button onClick={() => setShowBetModal(false)}>Cancel</button>
            </div>
          </div>
        )}

        <div className="fixed bottom-0 left-0 right-0">
          <div className="navbar-container">
            <NavItem
              icon={<Gamepad2 size={20} />}
              label="Matches"
              active={activeTab === "games"}
              onClick={() => setActiveTab("games")}
              to="/games"
            />
            <NavItem
              icon={<Receipt size={20} />}
              label="Bets"
              active={activeTab === "bets"}
              onClick={() => setActiveTab("bets")}
              to="/bets"
            />
            <NavItem
              icon={<UserCircle size={20} />}
              label="Profile"
              active={activeTab === "profile"}
              onClick={() => setActiveTab("profile")}
              to="/profile"
            />
          </div>
        </div>
      </div>
    </>
  );
}

function NavItem({ icon, label, active, onClick, to }) {
  return (
    <div className="navbar-item">
      <Link
        to={to}
        className={`nav-link ${active ? "active" : ""}`}
        onClick={onClick}
      >
        <div className="icon-container">{icon}</div>
        <span className="label">{label}</span>
      </Link>
    </div>
  );
}

export default Games;