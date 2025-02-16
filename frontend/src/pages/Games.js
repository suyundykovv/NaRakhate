"use client"

import { useState } from "react"
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
} from "lucide-react"
import { Link } from "react-router-dom"
import "./style_games.css"

function Games() {
  const [activeTab, setActiveTab] = useState("games")
  const [selectedLeague, setSelectedLeague] = useState(null)
  const [error, setError] = useState({ message: "", show: false })
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
  ]

  const leagueMatches = {
    "premier-league": [
      {
        time: "14:00, 12 August",
        team1: {
          name: "Manchester City",
          odds: "1.50",
          logo: "https://upload.wikimedia.org/wikipedia/en/thumb/e/eb/Manchester_City_FC_badge.svg/180px-Manchester_City_FC_badge.svg.png",
        },
        team2: {
          name: "Arsenal",
          odds: "2.75",
          logo: "https://upload.wikimedia.org/wikipedia/en/thumb/5/53/Arsenal_FC.svg/1200px-Arsenal_FC.svg.png",
        },
      },
    ],
    laliga: [
      {
        time: "16:15, 13 August",
        team1: {
          name: "Real Madrid",
          odds: "1.80",
          logo: "https://e7.pngegg.com/pngimages/552/638/png-clipart-real-madrid-c-f-la-liga-logo-real-madrid-logo-crown-logo-miscellaneous-sport-thumbnail.png",
        },
        team2: {
          name: "Barcelona",
          odds: "2.30",
          logo: "https://e7.pngegg.com/pngimages/836/292/png-clipart-fc-barcelona-fc-barcelona.png",
        },
      },
    ],
    bundesliga: [
      {
        time: "19:30, 17 August",
        team1: {
          name: "Bayern Munich",
          odds: "1.75",
          logo: "https://upload.wikimedia.org/wikipedia/commons/thumb/1/1b/FC_Bayern_M%C3%BCnchen_logo_%282017%29.svg/1200px-FC_Bayern_M%C3%BCnchen_logo_%282017%29.svg.png",
        },
        team2: {
          name: "Borussia Dortmund",
          odds: "2.20",
          logo: "https://upload.wikimedia.org/wikipedia/commons/thumb/6/67/Borussia_Dortmund_logo.svg/1200px-Borussia_Dortmund_logo.svg.png",
        },
      },
    ],
    "serie-a": [
      {
        time: "20:45, 18 August",
        team1: {
          name: "Juventus",
          odds: "1.90",
          logo: "https://e7.pngegg.com/pngimages/177/345/png-clipart-juventus-logo-juventus-f-c-serie-a-juventus-stadium-football-uefa-champions-league-football-text-sport-thumbnail.png",
        },
        team2: {
          name: "AC Milan",
          odds: "2.15",
          logo: "https://upload.wikimedia.org/wikipedia/commons/thumb/d/d0/Logo_of_AC_Milan.svg/1200px-Logo_of_AC_Milan.svg.png",
        },
      },
    ],
    "ligue-1": [
      {
        time: "21:00, 19 August",
        team1: {
          name: "PSG",
          odds: "1.65",
          logo: "https://upload.wikimedia.org/wikipedia/en/thumb/a/a7/Paris_Saint-Germain_F.C..svg/1200px-Paris_Saint-Germain_F.C..svg.png",
        },
        team2: {
          name: "Marseille",
          odds: "2.40",
          logo: "https://upload.wikimedia.org/wikipedia/commons/thumb/d/d8/Olympique_Marseille_logo.svg/1200px-Olympique_Marseille_logo.svg.png",
        },
      },
    ],
  }

  const handleOddsClick = () => {
    setError({ message: "You can't place bets until the card is linked", show: true })
    setTimeout(() => setError({ message: "", show: false }), 5000)
  }

  return (
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
            {leagueMatches[selectedLeague].map((match, index) => (
              <div key={index} className="match-item">
                <div className="match-header">
                  <span>{match.time}</span>
                  <div className="match-actions">
                    <ShoppingCart size={18} />
                    <Layout size={18} />
                  </div>
                </div>

                <div className="teams">
                  <div className="team team-left">
                    <img src={match.team1.logo || "/placeholder.svg"} alt={match.team1.name} className="team-logo" />
                    <span className="team-name">{match.team1.name}</span>
                  </div>
                  <div className="vs">VS</div>
                  <div className="team team-right">
                    <span className="team-name">{match.team2.name}</span>
                    <img src={match.team2.logo || "/placeholder.svg"} alt={match.team2.name} className="team-logo" />
                  </div>
                </div>

                <div className="odds">
                  <button className="odds-button" onClick={handleOddsClick}>
                    {match.team1.name} ({match.team1.odds})
                  </button>
                  <button className="odds-button" onClick={handleOddsClick}>
                    {match.team2.name} ({match.team2.odds})
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
  )
}

function NavItem({ icon, label, active, onClick, to }) {
  return (
    <div className="navbar-item">
      <Link to={to} className={`nav-link ${active ? "active" : ""}`} onClick={onClick}>
        <div className="icon-container">{icon}</div>
        <span className="label">{label}</span>
      </Link>
    </div>
  )
}

export default Games

