import React from "react";
import { Link } from "react-router-dom";

export default function NavItem({ icon, label, active, onClick, to }) {
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