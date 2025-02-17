import React from "react";
import { X } from "lucide-react";

export default function InfoPanel({ show, onClose }) {
  if (!show) return null;

  return (
    <div className="info-panel-overlay">
      <div className="info-panel">
        <div className="info-panel-header">
          <h3>Information</h3>
          <button 
            className="close-button"
            onClick={onClose}
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
  );
}