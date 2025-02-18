import React from 'react';
import { X } from 'lucide-react';
import './CharacterPersonalInfo.css';

const CharacterInfoPanel = ({ show, onClose, id, email, balance }) => {
  if (!show) return null;

  return (
    <div className="character-info-overlay">
      <div className="character-info-modal">
        <div className="character-info-header">
          <h2>Character Information</h2>
          <button className="close-button" onClick={onClose}>
            <X size={24} />
          </button>
        </div>
        <div className="character-info-content">
          <div className="info-item">
            <span className="info-label">Login ID:</span>
            <span className="info-value">{id}</span>
          </div>
          <div className="info-item">
            <span className="info-label">Email:</span>
            <span className="info-value">{email}</span>
          </div>
          <div className="info-item">
            <span className="info-label">Balance:</span>
            <span className="info-value">${balance.toFixed(2)}</span>
          </div>
        </div>
      </div>
    </div>
  );
};

export default CharacterInfoPanel;
