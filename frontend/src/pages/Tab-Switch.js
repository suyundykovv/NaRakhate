export function TabSwitch({ activeTab, onChange }) {
    return (
      <div className="mx-4 bg-white rounded-full p-1 flex">
        <button
          className={`flex-1 py-2 px-4 rounded-full text-sm transition-colors ${
            activeTab === "slip" ? "bg-[#B4E4E4]" : ""
          }`}
          onClick={() => onChange("slip")}
        >
          Bets Slip
        </button>
        <button
          className={`flex-1 py-2 px-4 rounded-full text-sm transition-colors ${
            activeTab === "history" ? "bg-[#B4E4E4]" : ""
          }`}
          onClick={() => onChange("history")}
        >
          Bet History
        </button>
      </div>
    );
  }
  