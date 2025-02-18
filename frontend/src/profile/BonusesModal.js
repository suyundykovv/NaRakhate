import { X } from "lucide-react"
import "./BonusModal.css"

const BonusesModal = ({ show, onClose }) => {
  if (!show) return null

  return (
    <div className="bonuses-overlay">
      <div className="bonuses-modal">
        <div className="bonuses-header">
          <h2>Bonuses</h2>
          <button className="close-button" onClick={onClose}>
            <X size={24} />
          </button>
        </div>
        <div className="bonuses-content">
          <p>There are currently no promotions or contests available.</p>
          <p>Check back later for exciting offers!</p>
        </div>
      </div>
    </div>
  )
}

export default BonusesModal

