import React from "react";
import { X, Check } from "lucide-react";

export default function WithdrawModal({ show, onClose, cards, selectedCard, setSelectedCard, amount, setAmount, success, setSuccess }) {
  const handleWithdrawSubmit = (e) => {
    e.preventDefault();
    setSuccess("Amount was successfully transferred to the card");
    setTimeout(() => {
      setSuccess("");
      onClose();
    }, 3000);
  };

  if (!show) return null;

  return (
    <div className="info-panel-overlay">
      <div className="info-panel">
        <div className="info-panel-header">
          <h3>Withdraw</h3>
          <button 
            className="close-button"
            onClick={() => {
              onClose();
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
  );
}