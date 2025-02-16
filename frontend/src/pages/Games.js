import React, { useState, useEffect } from "react";
import {
  ArrowLeft,
  ShoppingCart,
  Layout,
  ChevronRight,
  Gamepad2,
  Receipt,
  UserCircle,
  AlertCircle,
  X,
} from "lucide-react";
import { Link } from "react-router-dom";
import "./style_games.css";
import { fetchMatchesByLeague } from '../api'; // Импортируем функцию для получения матчей

function Games() {
  const [activeTab, setActiveTab] = useState("games");
  const [selectedLeague, setSelectedLeague] = useState(null);
  const [error, setError] = useState({ message: "", show: false });
  const [leagueMatches, setLeagueMatches] = useState({}); // Состояние для хранения данных о матчах

  const leagues = [
    {
      id: "premier-league",
      name: "Premier League",
      country: "England",
      logo: "https://w7.pngwing.com/pngs/157/277/png-transparent-lion-king-illustration-2016u201317-premier-league-1999u20132000-fa-premier-league-2017u201318-premier-league-english-football-league-chelsea-f-c-premier-league-file-purple-violet-logo-thumbnail.png",
    },
    {
      id: "laliga",
      name: "La Liga",
      country: "Spain",
      logo: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRVrpWnhRW2DNpZXEbhP-qSnNPRGJdcv0cGxA&s",
    },
    {
      id: "bundesliga",
      name: "Bundesliga",
      country: "Germany",
      logo: "https://upload.wikimedia.org/wikipedia/en/thumb/d/df/Bundesliga_logo_%282017%29.svg/1200px-Bundesliga_logo_%282017%29.svg.png",
    },
    {
      id: "serie-a",
      name: "Serie A",
      country: "Italy",
      logo: "https://img2.freepng.ru/20180811/tqw/d25f3cd241a52cb74c6a8bc79f48402d.webp",
    },
    {
      id: "ligue-1",
      name: "Ligue 1",
      country: "France",
      logo: "https://upload.wikimedia.org/wikipedia/commons/thumb/5/5e/Ligue1.svg/1200px-Ligue1.svg.png",
    },
  ];

  useEffect(() => {
    if (selectedLeague) {
      const loadMatches = async () => {
        try {
          const matches = await fetchMatchesByLeague(selectedLeague);
          setLeagueMatches({ [selectedLeague]: matches });
        } catch (error) {
          console.error('Failed to load matches:', error);
          setError({ message: "Failed to load matches", show: true });
        }
      };

      loadMatches();
    }
  }, [selectedLeague]);

 const handleOddsClick = () => {
   setError({ message: "You can't place bets until the card is linked", show: true });
   setTimeout(() => setError({ message: "", show: false }), 5000);
 };

  return (
    <>
      <div className="container">
        {error.show && (
          <div className={`error-notification ${error.show ? "show" : ""}`}>
            <AlertCircle className="error-notification-icon" size={20} />
            <span className="error-notification-message">{error.message}</span>
            <button className="error-notification-close" onClick={() => setError({ ...error, show: false })}>
              <X size={16} />
            </button>
          </div>
        )}

        <div className="header">
          {selectedLeague ? (
            <button onClick={() => setSelectedLeague(null)} className="back-button">
              <ArrowLeft className="text-gray-600" />
            </button>
          ) : null}
          <h1>{selectedLeague ? leagues.find((l) => l.id === selectedLeague)?.name : "Leagues"}</h1>
        </div>

        <div className="content">
          {!selectedLeague ? (
            <div className="league-list">
              {leagues.map((league) => (
                <button key={league.id} className="league-item" onClick={() => setSelectedLeague(league.id)}>
                  <div className="info">
                    <img src={league.logo || "/placeholder.svg"} alt={league.name} className="league-logo" />
                    <div className="league-text">
                      <span className="league-name">{league.name}</span>
                      <span className="league-country">{league.country}</span>
                    </div>
                  </div>
                  <ChevronRight className="text-gray-400" />
                </button>
              ))}
            </div>
          ) : (
            <div className="match-list">
              {leagueMatches[selectedLeague]?.map((match, index) => (
                <div key={index} className="match-item">
                  <div className="match-header">
                    <span>{new Date(match.FixtureDetails.Date).toLocaleString()}</span>
                    <div className="match-actions">
                      <ShoppingCart size={18} />
                      <Layout size={18} />
                    </div>
                  </div>

                  <div className="teams">
                    <div className="team team-left">
                      <img src={match.Teams.Home.Logo} alt={match.Teams.Home.Name} className="team-logo" />
                      <span className="team-name">{match.Teams.Home.Name}</span>
                    </div>
                    <div className="vs">VS</div>
                    <div className="team team-right">
                      <span className="team-name">{match.Teams.Away.Name}</span>
                      <img src={match.Teams.Away.Logo} alt={match.Teams.Away.Name} className="team-logo" />
                    </div>
                  </div>

                  <div className="odds">
                    <button className="odds-button" onClick={handleOddsClick}>
                      {match.Teams.Home.Name} ({match.Odd.HomeWin})
                    </button>
                    <button className="odds-button" onClick={handleOddsClick}>
                      {match.Teams.Away.Name} ({match.Odd.AwayWin})
                    </button>
                  </div>
                </div>
              ))}
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
    </>
  );
}

function NavItem({ icon, label, active, onClick, to }) {
  return (
    <div className="navbar-item">
      <Link to={to} className={`nav-link ${active ? "active" : ""}`} onClick={onClick}>
        <div className="icon-container">{icon}</div>
        <span className="label">{label}</span>
      </Link>
    </div>
  );
}

export default Games;
