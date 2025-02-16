import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from "react-router-dom";
import { Gamepad2, Receipt, UserCircle, Clock, Timer, FolderRoot as Football } from 'lucide-react';
import './style.css';
import { fetchBets } from './api'; // Импортируем функцию для получения ставок

function App() {
  const [activeTab, setActiveTab] = useState('bets');
  const navigate = useNavigate();
  const [matchMinutes, setMatchMinutes] = useState({});
  const [betHistory, setBetHistory] = useState([]); // Состояние для хранения данных о ставках

  // Загрузка данных о ставках при монтировании компонента
  useEffect(() => {
    const loadBets = async () => {
      try {
        const bets = await fetchBets();
        setBetHistory(bets);
      } catch (error) {
        console.error('Failed to load bets:', error);
      }
    };

    loadBets();
  }, []);
  
  const betHistory = [
    {
      id: 1,
      date: '13 August',
      homeTeam: 'Manchester United',
      awayTeam: 'Liverpool',
      score: '2:1',
      coefficient: 1.85,
      status: 'finished',
      won: true,
      betAmount: 100,
      selectedTeam: 'Manchester United',
      winAmount: 185
    },
    {
      id: 2,
      date: '8 August',
      homeTeam: 'Arsenal',
      awayTeam: 'Chelsea',
      score: '0:2',
      coefficient: 2.1,
      status: 'live',
      betAmount: 50,
      selectedTeam: 'Arsenal',
      possibleWin: 105,
      currentCoefficient: 2.3,
      startTime: new Date(Date.now() - 45 * 60 * 1000) // Match started 45 minutes ago
    },
    {
      id: 3,
      date: '1 August',
      homeTeam: 'Barcelona',
      awayTeam: 'Real Madrid',
      score: 'vs',
      coefficient: 1.95,
      status: 'waiting',
      betAmount: 75,
      selectedTeam: 'Barcelona',
      possibleWin: 146.25
    },
    {
      id: 4,
      date: '3 July',
      homeTeam: 'Bayern Munich',
      awayTeam: 'Dortmund',
      score: '3:1',
      coefficient: 1.75,
      status: 'finished',
      won: false,
      betAmount: 200,
      selectedTeam: 'Dortmund',
      lostAmount: 200
    },
  ];

  useEffect(() => {
    // Update match minutes for live matches
    const interval = setInterval(() => {
      const newMatchMinutes = {};
      betHistory.forEach(bet => {
        if (bet.status === 'live' && bet.startTime) {
          const elapsedMinutes = Math.floor((Date.now() - new Date(bet.startTime).getTime()) / (60 * 1000));
          newMatchMinutes[bet.id] = elapsedMinutes;
        }
      });
      setMatchMinutes(newMatchMinutes);
    }, 1000);

    return () => clearInterval(interval);
  }, []);


  const handlePlaceBet = async (betData) => {
    navigate('/games');
    try {
      const newBet = await createBet(betData);
      setBetHistory([...betHistory, newBet]); // Обновляем состояние с новой ставкой
    } catch (error) {
      console.error('Failed to place bet:', error);
    }
  };

  const getStatusIcon = (status, betId) => {
    switch (status) {
      case 'waiting':
        return <Clock className="status-icon waiting" />;
      case 'live':
        return (
          <div className="live-indicator">
            <Football className="status-icon live" />
            <span className="match-minute">{matchMinutes[betId]}′</span>
          </div>
        );
      case 'finished':
        return null;
      default:
        return null;
    }
  };

  const renderBetInfo = (bet) => {
    switch (bet.status) {
      case 'finished':
        if (bet.won) {
          return (
            <div className="bet-info won">
              <div className="bet-amount">Bet: ${bet.betAmount}</div>
              <div className="win-amount">Won: ${bet.winAmount}</div>
              <div className="profit">Profit: ${bet.winAmount - bet.betAmount}</div>
            </div>
          );
        } else {
          return (
            <div className="bet-info lost">
              <div className="bet-amount">Bet: ${bet.betAmount}</div>
              <div className="lost-amount">Lost: ${bet.lostAmount}</div>
            </div>
          );
        }
      case 'live':
        return (
          <div className="bet-info live">
            <div className="bet-amount">Bet: ${bet.betAmount}</div>
            <div className="coefficient-change">
              <span className="old-coefficient">{bet.coefficient}</span>
              <span className="arrow">→</span>
              <span className="new-coefficient">{bet.currentCoefficient}</span>
            </div>
            <div className="possible-win">Possible win: ${bet.possibleWin}</div>
          </div>
        );
      case 'waiting':
        return (
          <div className="bet-info waiting">
            <div className="bet-amount">Bet: ${bet.betAmount}</div>
            <div className="coefficient">Coefficient: {bet.coefficient}</div>
            <div className="possible-win">Possible win: ${bet.possibleWin}</div>
          </div>
        );
      default:
        return null;
    }
  };

  return (
    <div className="app-container">
      <div className="card">
        <div className="header">
          <h1 className="header-title">My Bets</h1>
        </div>

        <div className="tabs">
          <div className="tab-container">
            <button
              className={`tab ${activeTab === 'slip' ? 'active' : ''}`}
              onClick={() => setActiveTab('slip')}
            >
              Bets Slip
            </button>
            <button
              className={`tab ${activeTab === 'history' ? 'active' : ''}`}
              onClick={() => setActiveTab('history')}
            >
              Bet History
            </button>
          </div>
        </div>

        <div className="content">
          {activeTab === 'slip' ? (
            <div className="empty-state">
              <img
                src="https://images.unsplash.com/photo-1560421683-6856ea585c78?auto=format&fit=crop&w=300&q=80"
                alt="Empty state"
                className="empty-state-img"
              />
              <p className="empty-state-text">
                You currently do not have any booked games
              </p>
              <button className="place-bet-btn" onClick={handlePlaceBet}>
                Place a bet
              </button>
            </div>
          ) : (
            <div className="bet-history">
              {betHistory.map((bet) => (
                <div key={bet.id} className={`bet-item ${bet.status}`}>
                  <div className="bet-header">
                    <span className="bet-date">{new Date(bet.createdAt).toLocaleDateString()}</span>
                    <div className="bet-status-indicator">
                      {getStatusIcon(bet.status, bet.id)}
                      <span className="coefficient">{bet.oddValue}</span>
                    </div>
                  </div>

                  <div className="bet-content">
                    <div className="match-info">
                      <div className="teams">
                        <span className={`team ${bet.selectedTeam === bet.homeTeam ? 'selected' : ''}`}>
                          {bet.homeTeam}
                        </span>
                        <span className="score">{bet.score}</span>
                        <span className={`team ${bet.selectedTeam === bet.awayTeam ? 'selected' : ''}`}>
                          {bet.awayTeam}
                        </span>
                      </div>
                    </div>
                    {renderBetInfo(bet)}
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
    </div>
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
        <div className="icon-container">
          {icon}
        </div>
        <span className="label">{label}</span>
      </Link>
    </div>
  );
}

export default App;