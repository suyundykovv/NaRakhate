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
import { createBet } from "../api"; // Импортируем функцию для создания ставки

function Games() {
  const [activeTab, setActiveTab] = useState("games");
  const [selectedLeague, setSelectedLeague] = useState(null);
  const [events, setEvents] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [selectedBet, setSelectedBet] = useState(null); // Выбранная ставка (матч и тип ставки)
  const [betAmount, setBetAmount] = useState(""); // Сумма ставки
  const [showBetModal, setShowBetModal] = useState(false); // Показ модального окна

  const fetchEvents = async () => {
    try {
      const response = await fetch('http://127.0.0.1:8080/getEvents');
      if (!response.ok) throw new Error('Ошибка загрузки данных');
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
    const loadEvents = async () => {
      const eventsData = await fetchEvents();
      setEvents(eventsData);
    };
    loadEvents();
  }, []);

  // Обработчик клика на кнопку ставки
  const handleBetClick = (event, oddSelection, oddValue) => {
    setSelectedBet({ event, oddSelection, oddValue });
    setShowBetModal(true);
  };

  // Обработчик отправки ставки
  const handlePlaceBet = async () => {
    if (!selectedBet || !betAmount || isNaN(betAmount) || betAmount <= 0) {
      alert("Введите корректную сумму ставки");
      return;
    }

    const betData = {
      event_id: selectedBet.event.id,
      odd_selection: selectedBet.oddSelection,
      odd_value: selectedBet.oddValue,
      amount: parseFloat(betAmount),
    };

    try {
      await createBet(betData);
      alert("Ставка успешно создана!");
      setShowBetModal(false);
      setSelectedBet(null);
      setBetAmount("");
    } catch (error) {
      console.error("Ошибка при создании ставки:", error);
      alert("Не удалось создать ставку");
    }
  };

  return (
    <>
      <div className="container">
        <div className="header">
          <h1>Футбольные матчи</h1>
        </div>

        <div className="content">
          {loading && <div className="loading">Загрузка матчей...</div>}
          {error && <div className="error">Ошибка: {error}</div>}
          {!loading && !error && events.length === 0 && (
            <div className="empty">Нет доступных матчей</div>
          )}

          <div className="match-list">
            {events.map((event) => (
              <div key={event.id} className="match-item">
                <div className="match-header">
                  <span>
                    {new Date(event.start_time).toLocaleDateString('ru-RU', {
                      day: 'numeric',
                      month: 'long',
                      hour: '2-digit',
                      minute: '2-digit'
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
                      {event.name.split(' vs ')[0]}
                    </span>
                  </div>
                  <div className="vs">VS</div>
                  <div className="team team-right">
                    <span className="team-name">
                      {event.name.split(' vs ')[1]}
                    </span>
                  </div>
                </div>

                <div className="betting-buttons">
                  <button
                    className="bet-button"
                    onClick={() => handleBetClick(event, "home", event.odds.home_win)}
                  >
                    П1 ({event.odds.home_win.toFixed(2)})
                  </button>
                  <button
                    className="bet-button"
                    onClick={() => handleBetClick(event, "draw", event.odds.draw)}
                  >
                    Н ({event.odds.draw.toFixed(2)})
                  </button>
                  <button
                    className="bet-button"
                    onClick={() => handleBetClick(event, "away", event.odds.away_win)}
                  >
                    П2 ({event.odds.away_win.toFixed(2)})
                  </button>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Модальное окно для ставки */}
        {showBetModal && (
          <div className="bet-modal">
            <div className="bet-modal-content">
              <h3>Сделать ставку</h3>
              <p>Матч: {selectedBet.event.name}</p>
              <p>Ставка: {selectedBet.oddSelection === "home" ? "П1" : selectedBet.oddSelection === "away" ? "П2" : "Н"}</p>
              <p>Коэффициент: {selectedBet.oddValue.toFixed(2)}</p>
              <input
                type="number"
                placeholder="Сумма ставки"
                value={betAmount}
                onChange={(e) => setBetAmount(e.target.value)}
              />
              <button onClick={handlePlaceBet}>Сделать ставку</button>
              <button onClick={() => setShowBetModal(false)}>Отмена</button>
            </div>
          </div>
        )}

        <div className="fixed bottom-0 left-0 right-0">
          <div className="navbar-container">
            <NavItem
              icon={<Gamepad2 size={20} />}
              label="Матчи"
              active={activeTab === "games"}
              onClick={() => setActiveTab("games")}
              to="/games"
            />
            <NavItem
              icon={<Receipt size={20} />}
              label="Ставки"
              active={activeTab === "bets"}
              onClick={() => setActiveTab("bets")}
              to="/bets"
            />
            <NavItem
              icon={<UserCircle size={20} />}
              label="Профиль"
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