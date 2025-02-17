import React from "react";

export default function BalanceCard({ label, amount }) {
  return (
    <div className="balance-card">
      <p className="balance-label">{label}</p>
      <p className="balance-amount">{amount}</p>
    </div>
  );
}