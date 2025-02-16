import React, { useState } from "react";
import { Link } from "react-router-dom";
import { ChevronRight, Gift, Headphones, CreditCard, User, Gamepad2, Receipt, UserCircle, X, Camera, Check, Trash2, Plus, AlertCircle } from "lucide-react";
import "./styles.css";

export default function Profile() {
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
    firstName: "",
    lastName: "",
    avatarUrl: ""
  });

  const [cardData, setCardData] = useState({
    cardNumber: "",
    cardHolder: "",
    expiryDate: "",
    cvv: ""
  });

  const handleCardNumberChange = (e) => {
    let value = e.target.value.replace(/\D/g, "");
    value = value.replace(/(\d{4})/g, "$1 ").trim();
    if (value.length <= 19) {
      setCardData(prev => ({ ...prev, cardNumber: value }));
    }
  };

  const handleExpiryDateChange = (e) => {
    let value = e.target.value.replace(/\D/g, "");
    if (value.length > 2) {
      value = value.slice(0, 2) + "/" + value.slice(2, 4);
    }
    if (value.length <= 5) {
      setCardData(prev => ({ ...prev, expiryDate: value }));
    }
  };

  const handleCvvChange = (e) => {
    let value = e.target.value.replace(/\D/g, "");
    if (value.length <= 3) {
      setCardData(prev => ({ ...prev, cvv: value }));
    }
  };

  const handleCardSubmit = (e, type) => {
    e.preventDefault();
    
    if (cards.length >= 3) {
      setError("Maximum of 3 cards allowed");
      return;
    }

    if (cards.some(card => card.cardNumber === cardData.cardNumber.replace(/\s/g, ''))) {
      setError("This card is already added");
      return;
    }

    const newCard = {
      id: Math.random().toString(36).substr(2, 9),
      type,
      ...cardData
    };

    setCards([...cards, newCard]);
    setCardData({
      cardNumber: "",
      cardHolder: "",
      expiryDate: "",
      cvv: ""
    });
    setSelectedCard(null);
    setShowCardManagement(false);
    setError("");
  };

  const handleDeleteCard = (id) => {
    setCards(cards.filter(card => card.id !== id));
  };

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

  const handleDepositSubmit = (e) => {
    e.preventDefault();
    setError("Insufficient funds on the card");
    setTimeout(() => setError(""), 3000);
  };

  const handleWithdrawSubmit = (e) => {
    e.preventDefault();
    setSuccess("Amount was successfully transferred to the card");
    setTimeout(() => {
      setSuccess("");
      setShowWithdrawModal(false);
    }, 3000);
  };

  const handleAvatarChange = (e) => {
    const file = e.target.files?.[0];
    if (file) {
      const imageUrl = URL.createObjectURL(file);
      setProfileData(prev => ({
        ...prev,
        avatarUrl: imageUrl
      }));
    }
  };

  const handlePersonalInfoSubmit = (e) => {
    e.preventDefault();
    setShowPersonalInfo(false);
  };

  return (
    <div className="container">
      {/* Profile Section */}
      <div className="profile-header">
        <div className="profile-info">
          <div className="profile-img-container">
            <img
              src={profileData.avatarUrl || "https://via.placeholder.com/100x100"}
              alt="Profile"
              className="profile-img"
            />
          </div>
          <div>
            <h2 className="profile-name">
              {profileData.firstName && profileData.lastName 
                ? `${profileData.firstName} ${profileData.lastName}`
                : "Add Your Name"}
            </h2>
            <p className="profile-id">ID 090335</p>
          </div>
        </div>

        {/* Balance Cards */}
        <div className="balance-container">
          <BalanceCard label="Available balance" amount="$ 355.00" />
          <BalanceCard label="Debt owed" amount="$ 5.00" />
          <BalanceCard label="Cum. trans." amount="$ 0.00" />
        </div>
      </div>

      {/* Main Content */}
      <div className="main-content">
        {/* Action Buttons */}
        <div className="buttons-container">
          <button className="action-btn" onClick={handleDeposit}>Deposit</button>
          <button className="action-btn" onClick={handleWithdraw}>Withdraw</button>
        </div>

        {/* Menu List */}
        <div className="menu-list">
          <MenuItem 
            icon={<User size={20} />} 
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
                  {cards.map(card => (
                    <div key={card.id} className={`saved-card ${card.type}`}>
                      <div className="saved-card-info">
                        <div className="card-brand-logo">{card.type.toUpperCase()}</div>
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
      {showPersonalInfo && (
        <div className="info-panel-overlay">
          <div className="info-panel">
            <div className="info-panel-header">
              <h3>Personal Information</h3>
              <button 
                className="close-button"
                onClick={() => setShowPersonalInfo(false)}
              >
                <X size={24} />
              </button>
            </div>
            <div className="info-panel-content">
              <form onSubmit={handlePersonalInfoSubmit}>
                <div className="avatar-upload">
                  <div className="avatar-preview">
                    <img
                      src={profileData.avatarUrl || "https://via.placeholder.com/100x100"}
                      alt="Avatar preview"
                    />
                    <label className="avatar-upload-button">
                      <Camera size={20} />
                      <input
                        type="file"
                        accept="image/*"
                        onChange={handleAvatarChange}
                        className="hidden"
                      />
                    </label>
                  </div>
                </div>

                <div className="form-group">
                  <label htmlFor="firstName">First Name</label>
                  <input
                    type="text"
                    id="firstName"
                    value={profileData.firstName}
                    onChange={(e) => setProfileData(prev => ({
                      ...prev,
                      firstName: e.target.value
                    }))}
                    className="form-input"
                    placeholder="Enter your first name"
                  />
                </div>

                <div className="form-group">
                  <label htmlFor="lastName">Last Name</label>
                  <input
                    type="text"
                    id="lastName"
                    value={profileData.lastName}
                    onChange={(e) => setProfileData(prev => ({
                      ...prev,
                      lastName: e.target.value
                    }))}
                    className="form-input"
                    placeholder="Enter your last name"
                  />
                </div>

                <button type="submit" className="submit-button">
                  Save Changes
                </button>
              </form>
            </div>
          </div>
        </div>
      )}

      {/* Card Management Panel */}
      {showCardManagement && (
        <div className="info-panel-overlay">
          <div className="info-panel">
            <div className="info-panel-header">
              <h3>Card Management</h3>
              <button 
                className="close-button"
                onClick={() => {
                  setShowCardManagement(false);
                  setSelectedCard(null);
                }}
              >
                <X size={24} />
              </button>
            </div>
            <div className="info-panel-content">
              {!selectedCard ? (
                <div className="card-selection">
                  <h4>Select Card Type</h4>
                  <div className="card-options">
                    <div 
                      className="card-option visa"
                      onClick={() => setSelectedCard('visa')}
                    >
                      <div className="card-brand-logo">VISA</div>
                    </div>
                    <div 
                      className="card-option mastercard"
                      onClick={() => setSelectedCard('mastercard')}
                    >
                      <div className="card-brand-logo">MASTERCARD</div>
                    </div>
                  </div>
                </div>
              ) : (
                <div className={`card-form ${selectedCard}`}>
                  <button 
                    className="back-button"
                    onClick={() => setSelectedCard(null)}
                  >
                    ← Back to card selection
                  </button>
                  
                  <div className="card-preview">
                    <div className="card-brand-logo">{selectedCard.toUpperCase()}</div>
                    <div className="card-number">
                      {cardData.cardNumber || "•••• •••• •••• ••••"}
                    </div>
                    <div className="card-details">
                      <div>
                        <div className="card-label">CARD HOLDER</div>
                        <div className="card-value">{cardData.cardHolder || "YOUR NAME"}</div>
                      </div>
                      <div>
                        <div className="card-label">EXPIRES</div>
                        <div className="card-value">{cardData.expiryDate || "MM/YY"}</div>
                      </div>
                    </div>
                  </div>

                  <form onSubmit={(e) => handleCardSubmit(e, selectedCard)}>
                    <div className="form-group">
                      <label>Card Number</label>
                      <input
                        type="text"
                        value={cardData.cardNumber}
                        onChange={handleCardNumberChange}
                        placeholder="1234 5678 9012 3456"
                        className="form-input"
                        maxLength="19"
                      />
                    </div>

                    <div className="form-group">
                      <label>Card Holder Name</label>
                      <input
                        type="text"
                        value={cardData.cardHolder}
                        onChange={(e) => setCardData(prev => ({ 
                          ...prev, 
                          cardHolder: e.target.value.toUpperCase()
                        }))}
                        placeholder="JOHN DOE"
                        className="form-input"
                      />
                    </div>

                    <div className="form-row">
                      <div className="form-group">
                        <label>Expiry Date</label>
                        <input
                          type="text"
                          value={cardData.expiryDate}
                          onChange={handleExpiryDateChange}
                          placeholder="MM/YY"
                          className="form-input"
                          maxLength="5"
                        />
                      </div>

                      <div className="form-group">
                        <label>CVV</label>
                        <input
                          type="text"
                          value={cardData.cvv}
                          onChange={handleCvvChange}
                          placeholder="123"
                          className="form-input"
                          maxLength="3"
                        />
                      </div>
                    </div>

                    {error && (
                      <div className="error-message">
                        <AlertCircle size={16} />
                        {error}
                      </div>
                    )}

                    <button type="submit" className={`submit-button ${selectedCard}`}>
                      Add Card
                    </button>
                  </form>
                </div>
              )}
            </div>
          </div>
        </div>
      )}

      {/* Deposit Modal */}
      {showDepositModal && (
        <div className="info-panel-overlay">
          <div className="info-panel">
            <div className="info-panel-header">
              <h3>Deposit</h3>
              <button 
                className="close-button"
                onClick={() => {
                  setShowDepositModal(false);
                  setAmount("");
                  setError("");
                }}
              >
                <X size={24} />
              </button>
            </div>
            <div className="info-panel-content">
              <form onSubmit={handleDepositSubmit}>
                <div className="form-group">
                  <label>Select Card</label>
                  <select 
                    className="form-input"
                    value={selectedCard || ""}
                    onChange={(e) => setSelectedCard(e.target.value)}
                  >
                    <option value="">Select a card</option>
                    {cards.map(card => (
                      <option key={card.id} value={card.id}>
                        {card.type.toUpperCase()} •••• {card.cardNumber.slice(-4)}
                      </option>
                    ))}
                  </select>
                </div>

                <div className="form-group">
                  <label>Amount</label>
                  <input
                    type="number"
                    className="form-input"
                    value={amount}
                    onChange={(e) => setAmount(e.target.value)}
                    placeholder="Enter amount"
                    min="1"
                  />
                </div>

                {error && (
                  <div className="error-message">
                    <AlertCircle size={16} />
                    {error}
                  </div>
                )}

                <button 
                  type="submit" 
                  className="submit-button"
                  disabled={!selectedCard || !amount}
                >
                  Deposit
                </button>
              </form>
            </div>
          </div>
        </div>
      )}

      {/* Withdraw Modal */}
      {showWithdrawModal && (
        <div className="info-panel-overlay">
          <div className="info-panel">
            <div className="info-panel-header">
              <h3>Withdraw</h3>
              <button 
                className="close-button"
                onClick={() => {
                  setShowWithdrawModal(false);
                  setAmount("");
                  setSuccess("");
                }}
              >
                <X size={24} />
              </button>
            </div>
            <div className="info-panel-content">
              <form onSubmit={handleWithdrawSubmit}>
                <div className="form-group">
                  <label>Select Card</label>
                  <select 
                    className="form-input"
                    value={selectedCard || ""}
                    onChange={(e) => setSelectedCard(e.target.value)}
                  >
                    <option value="">Select a card</option>
                    {cards.map(card => (
                      <option key={card.id} value={card.id}>
                        {card.type.toUpperCase()} •••• {card.cardNumber.slice(-4)}
                      </option>
                    ))}
                  </select>
                </div>

                <div className="form-group">
                  <label>Amount</label>
                  <input
                    type="number"
                    className="form-input"
                    value={amount}
                    onChange={(e) => setAmount(e.target.value)}
                    placeholder="Enter amount"
                    min="1"
                  />
                </div>

                {success && (
                  <div className="success-message">
                    <Check size={16} />
                    {success}
                  </div>
                )}

                <button 
                  type="submit" 
                  className="submit-button"
                  disabled={!selectedCard || !amount}
                >
                  Withdraw
                </button>
              </form>
            </div>
          </div>
        </div>
      )}

      {/* Info Panel */}
      {showInfo && (
        <div className="info-panel-overlay">
          <div className="info-panel">
            <div className="info-panel-header">
              <h3>Information</h3>
              <button 
                className="close-button"
                onClick={() => setShowInfo(false)}
              >
                <X size={24} />
              </button>
            </div>
            <div className="info-panel-content">
              <h4>About Our Service</h4>
              <p>Welcome to our platform! We're dedicated to providing you with the best gaming experience possible.</p>
              
              <h4>Contact Support</h4>
              <p>24/7 Customer Support: support@example.com</p>
              <p>Phone: +1 (555) 123-4567</p>
              
              <h4>Important Links</h4>
              <ul>
                <li>Terms of Service</li>
                <li>Privacy Policy</li>
                <li>Responsible Gaming</li>
                <li>FAQ</li>
              </ul>
            </div>
          </div>
        </div>
      )}

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

function BalanceCard(props) {
  return (
    <div className="balance-card">
      <p className="balance-label">{props.label}</p>
      <p className="balance-amount">{props.amount}</p>
    </div>
  );
}

function MenuItem({ icon, label, onClick }) {
  return (
    <div className="menu-item" onClick={onClick}>
      <div className="menu-icon-text">
        {icon}
        <span>{label}</span>
      </div>
      <ChevronRight size={20} className="chevron-icon" />
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