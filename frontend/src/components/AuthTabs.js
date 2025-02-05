import { useState } from "react"
import LoginForm from "./LoginForm"
import RegisterForm from "./RegisterForm"

const AuthTabs = () => {
  const [activeTab, setActiveTab] = useState("login")

  return (
    <div style={styles.container}>
      {/* Вкладки */}
      <div style={styles.tabs}>
        <button style={activeTab === "login" ? styles.activeTab : styles.tab} onClick={() => setActiveTab("login")}>
          Вход
        </button>
        <button
          style={activeTab === "register" ? styles.activeTab : styles.tab}
          onClick={() => setActiveTab("register")}
        >
          Регистрация
        </button>
      </div>

      {/* Форма входа */}
      {activeTab === "login" && <LoginForm />}

      {/* Форма регистрации */}
      {activeTab === "register" && <RegisterForm />}
    </div>
  )
}

const styles = {
  container: {
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
    justifyContent: "center",
    minHeight: "100vh",
    backgroundColor: "#B4E4E4",
    padding: "20px",
  },
  tabs: {
    display: "flex",
    gap: "32px",
    marginBottom: "20px",
    backgroundColor: "white",
    padding: "20px 32px 0",
    borderTopLeftRadius: "16px",
    borderTopRightRadius: "16px",
  },
  tab: {
    padding: "8px 0",
    border: "none",
    background: "none",
    cursor: "pointer",
    fontSize: "16px",
    color: "#666",
    position: "relative",
  },
  activeTab: {
    padding: "8px 0",
    border: "none",
    background: "none",
    cursor: "pointer",
    fontSize: "16px",
    color: "#000",
    position: "relative",
    "&::after": {
      content: '""',
      position: "absolute",
      bottom: "-2px",
      left: 0,
      right: 0,
      height: "2px",
      backgroundColor: "#FF4B55",
    },
  },
}

export default AuthTabs

