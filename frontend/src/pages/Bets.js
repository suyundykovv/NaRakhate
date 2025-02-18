import React, { useState, useEffect, useCallback } from "react";
import { Link, useNavigate } from "react-router-dom";
import { Gamepad2, Receipt, UserCircle, Clock, FolderRoot as Football } from "lucide-react";
import "./style.css";

// Utility function to fetch data with better error handling
const fetchData = async (url, options = {}) => {
  try {
    const response = await fetch(url, options);
    if (!response.ok) {
      throw new Error(`Failed to fetch data from ${url}: ${response.statusText}`);
    }
    return await response.json();
  } catch (error) {
    console.error("Error fetching data:", error);
    throw error; // Re-throw to handle errors in calling functions
  }
};

// Fetch user ID
const fetchUserId = async () => {
  const token = localStorage.getItem("token");
  if (!token) throw new Error("No token found. Please log in again.");

  const userData = await fetchData("http://127.0.0.1:8080/getUser", {
    method: "GET",
    headers: { Authorization: `Bearer ${token}` },
  });
  return userData?.id;
};

// Fetch bets
const fetchBets = async () => {
  const allBets = await fetchData("http://127.0.0.1:8080/getBets");
  return Array.isArray(allBets) ? allBets : [];
};

// Fetch event data
const fetchEventData = async (eventId) => {
  return await fetchData(`http://127.0.0.1:8080/getEvent/${eventId}`);
};

// BetInfo Component
const BetInfo = React.memo(({ bet, event }) => {
  const homeTeam = event?.homeTeam || "Home Team";
  const awayTeam = event?.awayTeam || "Away Team";
  const homeGoals = event?.homeGoals || 0;
  const awayGoals = event?.awayGoals || 0;

  switch (bet.status) {
    case "finished":
      return (
        <div className={`bet-info ${bet.income > bet.amount ? "won" : "lost"}`}>
          <div className="bet-amount">Bet: ${bet.amount.toFixed(2)}</div>
          {bet.income > bet.amount ? (
            <>
              <div className="win-amount">Won: ${bet.income.toFixed(2)}</div>
              <div className="profit">Profit: ${(bet.income - bet.amount).toFixed(2)}</div>
            </>
          ) : (
            <div className="lost-amount">Lost: ${bet.amount.toFixed(2)}</div>
          )}
        </div>
      );
    case "live":
      return (
        <div className="bet-info live">
          <div className="bet-amount">Bet: ${bet.amount.toFixed(2)}</div>
          <div className="coefficient-change">
            <span className="old-coefficient">{bet.odd_value.toFixed(2)}</span>
            <span className="arrow">→</span>
            <span className="new-coefficient">{event?.odds[bet.odd_selection]?.toFixed(2) || "N/A"}</span>
          </div>
          <div className="possible-win">Possible win: ${bet.income.toFixed(2)}</div>
        </div>
      );
    case "waiting":
      return (
        <div className="bet-info waiting">
          <div className="bet-amount">Bet: ${bet.amount.toFixed(2)}</div>
          <div className="coefficient">Coefficient: {bet.odd_value.toFixed(2)}</div>
          <div className="possible-win">Possible win: ${bet.income.toFixed(2)}</div>
        </div>
      );
    default:
      return null;
  }
});

function App() {
  const [activeTab, setActiveTab] = useState("bets");
  const navigate = useNavigate();
  const [matchMinutes, setMatchMinutes] = useState({});
  const [betHistory, setBetHistory] = useState([]);
  const [userId, setUserId] = useState(null);
  const [events, setEvents] = useState({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
// Load data on mount
useEffect(() => {
  const loadData = async () => {
    setLoading(true);
    setError(null);

    try {
      // Fetch the user ID from the token
      const id = await fetchUserId();
      setUserId(id);

      // Fetch all bets
      const allBets = await fetchBets();

      // Filter bets to include only those belonging to the logged-in user
      const userBets = allBets.filter((bet) => bet.user_id === id);

      // Fetch event data for each bet
      const eventsData = {};
      for (const bet of userBets) {
        const eventData = await fetchEventData(bet.event_id);
        if (eventData) eventsData[bet.event_id] = eventData;
      }

      // Update state with filtered bets and event data
      setEvents(eventsData);
      setBetHistory(userBets);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  loadData();
}, [userId]);

  // Update match minutes for live matches
  useEffect(() => {
    const interval = setInterval(() => {
      const newMatchMinutes = {};
      betHistory.forEach((bet) => {
        if (bet.status === "live" && bet.startTime) {
          const elapsedMinutes = Math.floor(
            (Date.now() - new Date(bet.startTime).getTime()) / (60 * 1000)
          );
          newMatchMinutes[bet.id] = elapsedMinutes;
        }
      });
      setMatchMinutes(newMatchMinutes);
    }, 1000);

    return () => clearInterval(interval);
  }, [betHistory]);

  const getStatusIcon = useCallback((status, betId) => {
    switch (status) {
      case "waiting":
        return <Clock className="status-icon waiting" />;
      case "live":
        return (
          <div className="live-indicator">
            <Football className="status-icon live" />
            <span className="match-minute">{matchMinutes[betId]}′</span>
          </div>
        );
      default:
        return null;
    }
  }, [matchMinutes]);

  if (error) {
    return <div className="error">Error: {error}</div>;
  }

  return (
    <div className="app-container">
      <div className="card">
        <div className="header">
          <h1 className="header-title">My Bets</h1>
        </div>

        <div className="tabs">
          <div className="tab-container">
            <button
              className={`tab ${activeTab === "slip" ? "active" : ""}`}
              onClick={() => setActiveTab("slip")}
            >
              Bets Slip
            </button>
            <button
              className={`tab ${activeTab === "history" ? "active" : ""}`}
              onClick={() => setActiveTab("history")}
            >
              Bet History
            </button>
          </div>
        </div>

        <div className="content">
          {loading ? (
            <div className="loading">Loading bets...</div>
          ) : activeTab === "slip" ? (
            <div className="empty-state">
              <img
                src="https://images.unsplash.com/photo-1560421683-6856ea585c78?auto=format&fit=crop&w=300&q=80"
                alt="Empty state"
                className="empty-state-img"
              />
              <p className="empty-state-text">
                You currently do not have any booked games
              </p>
              <button className="place-bet-btn" onClick={() => navigate("/games")}>
                Place a bet
              </button>
            </div>
          ) : (
            <div className="bet-history">
              {betHistory.length === 0 ? (
                <div className="empty-bets">
                  <p>No bets found. Place a bet to get started!</p>
                  <button className="place-bet-btn" onClick={() => navigate("/games")}>
                    Place a bet
                  </button>
                </div>
              ) : (
                betHistory.map((bet) => {
                  const event = events[bet.event_id];
                  return (
                    <div key={bet.id} className={`bet-item ${bet.status}`}>
                      <div className="bet-header">
                        <span className="bet-date">
                          {new Date(bet.created_at).toLocaleDateString()}
                        </span>
                        <div className="bet-status-indicator">
                          {getStatusIcon(bet.status, bet.id)}
                          <span className="value">{" Income:" + bet.income.toFixed(2)}</span>
                        </div>
                      </div>
                  
                      <div className="bet-content">
                        <div className="match-info">
                          <div className="teams">
                            <div className="team team-left">
                              <span className="team-name">
                                {event?.name.split(" vs ")[0] || "Home Team"}
                              </span>
                              <span className="team-score">
                                { " _Goals:" + event?.home_goals || 0} {/* Fixed to `homeGoals` */}
                              </span>
                            </div>
                            <div className="status">{event?.status }</div>
                            <div className="team team-right">
                              <span className="team-name">
                                {event?.name.split(" vs ")[1] || "Away Team"}
                              </span>
                              <span className="team-score">
                                {" _Goals:" + event?.away_goals || 0} {/* Fixed to `awayGoals` */}
                              </span>
                            </div>
                          </div>
                        </div>
                        <BetInfo bet={bet} event={event} />
                      </div>
                    </div>
                  );
                })
              )}
            </div>
          )}
        </div>

        <div className="fixed bottom-0 left-0 right-0">
          <div className="navbar-container">
            <NavItem
              icon={<Gamepad2 size={20} />}
              label="Games"
              active={activeTab === "games"}
              onClick={() => setActiveTab("games")}
              to="/games"
            />
            <NavItem
              icon={<Receipt size={20} />}
              label="My Bets"
              active={activeTab === "bets"}
              onClick={() => setActiveTab("bets")}
              to="/bets"
            />
            <NavItem
              icon={<UserCircle size={20} />}
              label="My Profile"
              active={activeTab === "profile"}
              onClick={() => setActiveTab("profile")}
              to="/profile"
            />
          </div>
        </div>
      </div>
    </div>
  );
}

const NavItem = React.memo(({ icon, label, active, onClick, to }) => (
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
));

export default App;
