import React from "react";
import { X, AlertCircle } from "lucide-react";

export default function DepositModal({ show, onClose, cards, selectedCard, setSelectedCard, amount, setAmount, error, setError }) {
  const handleDepositSubmit = (e) => {
    e.preventDefault();
    setError("Insufficient funds on the card");
    setTimeout(() => setError(""), 3000);
  };

  if (!show) return null;

  return (
    <div className="info-panel-overlay">
      <div className="info-panel">
        <div className="info-panel-header">
          <h3>Deposit</h3>
          <button 
            className="close-button"
            onClick={() => {
              onClose();
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
  );
}