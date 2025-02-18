"use client"
import { Link, Route, Routes } from "react-router-dom"
import LoginForm from "./LoginForm"
import RegisterForm from "./RegisterForm"

const AuthTabs = ({ onLogin }) => {
  //const [activeTab, setActiveTab] = useState("login")

  return (
    <div style={styles.container}>
      <div style={styles.tabs}>
        <Link to="/login" style={styles.tab}>
          Вход
        </Link>
        <Link to="/register" style={styles.tab}>
          Регистрация
        </Link>
      </div>

      <Routes>
        <Route path="/login" element={<LoginForm onLogin={onLogin} />} />
        <Route path="/register" element={<RegisterForm />} />
      </Routes>
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
    textDecoration: "none",
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

