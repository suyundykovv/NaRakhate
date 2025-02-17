import React from "react";
import { ChevronRight } from "lucide-react";

export default function MenuItem({ icon, label, onClick }) {
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