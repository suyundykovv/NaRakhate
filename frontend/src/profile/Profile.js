import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import {
  Gamepad2,
  Receipt,
  UserCircle,
  Gift,
  CreditCard,
  Trash2,
  Plus,
  Headphones,
  X,
  Camera,
  Check,
  AlertCircle,
  LogOut,
} from "lucide-react"; // Import all required icons
import BalanceCard from "./BalanceCard";
import MenuItem from "./MenuItem";
import NavItem from "./NavItem";
import PersonalInfoPanel from "./PersonalInfoPanel";
import CardManagementPanel from "./CardManagementPanel";
import DepositModal from "./DepositModal";
import WithdrawModal from "./WithdrawModal";
import InfoPanel from "./InfoPanel";
import "../pages/styles.css"; // Adjust the path to your CSS file

export default function Profile({ onLogout }) {
  const [activeTab, setActiveTab] = useState("profile");
  const [showInfo, setShowInfo] = useState(false);
  const [showPersonalInfo, setShowPersonalInfo] = useState(false);
  const [showCardManagement, setShowCardManagement] = useState(false);
  const [showDepositModal, setShowDepositModal] = useState(false);
  const [showWithdrawModal, setShowWithdrawModal] = useState(false);
  const [selectedCard, setSelectedCard] = useState(null);
  const [amount, setAmount] = useState("");
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const [cards, setCards] = useState([]);
  const [profileData, setProfileData] = useState({
    id: "",
    username: "",
    email: "",
    role: "",
    cash: 0,
    avatarUrl: "https://via.placeholder.com/100x100", // Default avatar
  });

  // Fetch user data when the component mounts
  useEffect(() => {
    const fetchUserData = async () => {
      const token = localStorage.getItem("token"); // Get the token from localStorage
      if (!token) {
        setError("No token found. Please log in again.");
        return;
      }

      try {
        const response = await fetch("http://127.0.0.1:8080/getUser", {
          method: "GET",
          headers: {
            Authorization: `Bearer ${token}`, // Include the token in the request headers
          },
        });

        if (!response.ok) {
          throw new Error("Failed to fetch user data");
        }

        const userData = await response.json();
        setProfileData({
          id: userData.id,
          username: userData.username,
          email: userData.email,
          role: userData.role,
          cash: userData.cash,
          avatarUrl: "https://via.placeholder.com/100x100", // You can update this if the backend provides an avatar URL
        });
      } catch (err) {
        setError("Error fetching user data: " + err.message);
      }
    };

    fetchUserData();
  }, []); // Empty dependency array ensures this runs only once on mount

  const handleDeposit = () => {
    if (cards.length === 0) {
      setShowCardManagement(true);
      return;
    }
    setShowDepositModal(true);
  };

  const handleWithdraw = () => {
    if (cards.length === 0) {
      setShowCardManagement(true);
      return;
    }
    setShowWithdrawModal(true);
  };

  const handleDeleteCard = (id) => {
    setCards(cards.filter((card) => card.id !== id));
  };

  return (
    <div className="container">
      {/* Profile Section */}
      <div className="profile-header">
        <div className="profile-info">
          <div className="profile-img-container">
            <img
              src={profileData.avatarUrl}
              alt="Profile"
              className="profile-img"
            />
          </div>
          <div>
            <h2 className="profile-name">
              {profileData.username || "Add Your Name"}
            </h2>
            <p className="profile-id">ID: {profileData.id}</p>
            <p className="profile-email">Email: {profileData.email}</p>
            <p className="profile-role">Role: {profileData.role}</p>
          </div>
        </div>

        {/* Logout Button */}
        <button className="logout-button" onClick={onLogout}>
          <LogOut size={20} />
          Logout
        </button>

        {/* Balance Cards */}
        <div className="balance-container">
          <BalanceCard label="Available balance" amount={`$ ${profileData.cash.toFixed(2)}`} />
          <BalanceCard label="Debt owed" amount="$ 5.00" />
          <BalanceCard label="Cum. trans." amount="$ 0.00" />
        </div>
      </div>

      {/* Main Content */}
      <div className="main-content">
        {/* Action Buttons */}
        <div className="buttons-container">
          <button className="action-btn" onClick={handleDeposit}>
            Deposit
          </button>
          <button className="action-btn" onClick={handleWithdraw}>
            Withdraw
          </button>
        </div>

        {/* Menu List */}
        <div className="menu-list">
          <MenuItem
            icon={<UserCircle size={20} />}
            label="Personal Info"
            onClick={() => setShowPersonalInfo(true)}
          />
          <MenuItem icon={<Gift size={20} />} label="Bonuses" />
          <div className="card-management-section">
            <div className="card-management-header">
              <MenuItem
                icon={<CreditCard size={20} />}
                label="Card Management"
                onClick={() => setShowCardManagement(true)}
              />
              {cards.length > 0 && (
                <div className="cards-list">
                  {cards.map((card) => (
                    <div key={card.id} className={`saved-card ${card.type}`}>
                      <div className="saved-card-info">
                        <div className="card-brand-logo">
                          {card.type.toUpperCase()}
                        </div>
                        <div className="saved-card-number">
                          •••• {card.cardNumber.slice(-4)}
                        </div>
                      </div>
                      <button
                        className="delete-card-btn"
                        onClick={() => handleDeleteCard(card.id)}
                      >
                        <Trash2 size={16} />
                      </button>
                    </div>
                  ))}
                  {cards.length < 3 && (
                    <button
                      className="add-card-btn"
                      onClick={() => setShowCardManagement(true)}
                    >
                      <Plus size={20} />
                      Add Card
                    </button>
                  )}
                </div>
              )}
            </div>
          </div>
          <MenuItem
            icon={<Headphones size={20} />}
            label="Info"
            onClick={() => setShowInfo(true)}
          />
        </div>
      </div>

      {/* Personal Info Panel */}
      <PersonalInfoPanel
        show={showPersonalInfo}
        onClose={() => setShowPersonalInfo(false)}
        profileData={profileData}
        setProfileData={setProfileData}
      />

      {/* Card Management Panel */}
      <CardManagementPanel
        show={showCardManagement}
        onClose={() => setShowCardManagement(false)}
        cards={cards}
        setCards={setCards}
        selectedCard={selectedCard}
        setSelectedCard={setSelectedCard}
      />

      {/* Deposit Modal */}
      <DepositModal
        show={showDepositModal}
        onClose={() => setShowDepositModal(false)}
        cards={cards}
        selectedCard={selectedCard}
        setSelectedCard={setSelectedCard}
        amount={amount}
        setAmount={setAmount}
        error={error}
        setError={setError}
      />

      {/* Withdraw Modal */}
      <WithdrawModal
        show={showWithdrawModal}
        onClose={() => setShowWithdrawModal(false)}
        cards={cards}
        selectedCard={selectedCard}
        setSelectedCard={setSelectedCard}
        amount={amount}
        setAmount={setAmount}
        success={success}
        setSuccess={setSuccess}
      />

      {/* Info Panel */}
      <InfoPanel show={showInfo} onClose={() => setShowInfo(false)} />

      {/* Bottom Navigation */}
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
  );
}