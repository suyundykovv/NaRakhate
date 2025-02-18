"use client"

import { useState, useEffect } from "react"
import { Gamepad2, Receipt, UserCircle, Gift, Headphones, LogOut } from "lucide-react"
import BalanceCard from "./BalanceCard"
import MenuItem from "./MenuItem"
import NavItem from "./NavItem"
import PersonalInfoPanel from "./PersonalInfoPanel"
import BonusesModal from "./BonusesModal"
import InfoPanel from "./InfoPanel"
import "../pages/styles.css"

export default function Profile({ onLogout }) {
  const [activeTab, setActiveTab] = useState("profile")
  const [showInfo, setShowInfo] = useState(false)
  const [showPersonalInfo, setShowPersonalInfo] = useState(false)
  const [showBonuses, setShowBonuses] = useState(false)

  const [error, setError] = useState("")

  const [profileData, setProfileData] = useState({
    id: "",
    username: "",
    email: "",
    role: "",
    cash: 0,
  })

  useEffect(() => {
    const fetchUserData = async () => {
      const token = localStorage.getItem("token")
      if (!token) {
        setError("No token found. Please log in again.")
        return
      }

      try {
        const response = await fetch("http://127.0.0.1:8080/getUser", {
          method: "GET",
          headers: {
            Authorization: `Bearer ${token}`,
          },
        })

        if (!response.ok) {
          throw new Error("Failed to fetch user data")
        }

        const userData = await response.json()
        setProfileData({
          id: userData.id,
          username: userData.username,
          email: userData.email,
          role: userData.role,
          cash: userData.cash,
        })
      } catch (err) {
        setError("Error fetching user data: " + err.message)
      }
    }

    fetchUserData()
  }, [])

  return (
    <div className="container">
      <div className="profile-header">
        <div className="profile-info">
          <div>
            <h2 className="profile-name">{profileData.username || "Add Your Name"}</h2>
            <p className="profile-id">ID: {profileData.id}</p>
            <p className="profile-email">Email: {profileData.email}</p>
            <p className="profile-role">Role: {profileData.role}</p>
          </div>
        </div>

        <button className="logout-button" onClick={onLogout}>
          <LogOut size={20} />
          Logout
        </button>

        <div className="balance-container">
          <BalanceCard label="Available balance" amount={`$ ${profileData.cash.toFixed(2)}`} />
          <BalanceCard label="Debt owed" amount="$ 5.00" />
          <BalanceCard label="Cum. trans." amount="$ 0.00" />
        </div>
      </div>

      <div className="main-content">
        <div className="menu-list">
          <MenuItem icon={<UserCircle size={20} />} label="Personal Info" onClick={() => setShowPersonalInfo(true)} />
          <MenuItem icon={<Gift size={20} />} label="Bonuses" onClick={() => setShowBonuses(true)} />
          <MenuItem icon={<Headphones size={20} />} label="Info" onClick={() => setShowInfo(true)} />
        </div>
      </div>

      <PersonalInfoPanel
        show={showPersonalInfo}
        onClose={() => setShowPersonalInfo(false)}
        id={profileData.id}
        email={profileData.email}
        balance={profileData.cash}
      />

      <BonusesModal show={showBonuses} onClose={() => setShowBonuses(false)} />

      <InfoPanel show={showInfo} onClose={() => setShowInfo(false)} />

      <div className="fixed bottom-0 left-0 right-0">
        <div className="navbar-container">
          <NavItem
            icon={<Gamepad2 size={20} />}
            label="Games"
            active={activeTab === "games"}
            onClick={() => setActiveTab("games")}
            to="/games"
          />
          <NavItem
            icon={<Receipt size={20} />}
            label="My Bets"
            active={activeTab === "bets"}
            onClick={() => setActiveTab("bets")}
            to="/bets"
          />
          <NavItem
            icon={<UserCircle size={20} />}
            label="My Profile"
            active={activeTab === "profile"}
            onClick={() => setActiveTab("profile")}
            to="/profile"
          />
        </div>
      </div>
    </div>
  )
}

