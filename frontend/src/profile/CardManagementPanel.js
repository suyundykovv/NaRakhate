import React from "react";
import { X, Trash2, Plus, AlertCircle } from "lucide-react";

export default function CardManagementPanel({ show, onClose, cards, setCards, selectedCard, setSelectedCard }) {
  const [cardData, setCardData] = React.useState({
    cardNumber: "",
    cardHolder: "",
    expiryDate: "",
    cvv: ""
  });

  const [error, setError] = React.useState("");

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
    onClose();
    setError("");
  };

  const handleDeleteCard = (id) => {
    setCards(cards.filter(card => card.id !== id));
  };

  if (!show) return null;

  return (
    <div className="info-panel-overlay">
      <div className="info-panel">
        <div className="info-panel-header">
          <h3>Card Management</h3>
          <button 
            className="close-button"
            onClick={() => {
              onClose();
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
  );
}