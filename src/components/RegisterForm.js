import { useState } from "react"
import axios from "axios"  // Импортируем axios для отправки запросов

const LoginForm = () => {
  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")
  const [error, setError] = useState("")  // Для отображения ошибки

  const handleSubmit = async (e) => {
    e.preventDefault()
    console.log("Login:", username)
    console.log("Password:", password)

    try {
      // Отправляем POST запрос на сервер
      const response = await axios.post("http://localhost:8080/login", {
        username,
        password
      })

      // Сохраняем JWT токен в localStorage
      localStorage.setItem("token", response.data.token)

      // Перенаправляем или выводим успешное сообщение
      console.log("Login successful!", response.data.token)
    } catch (error) {
      setError("Invalid username or password")  // Показать ошибку, если она есть
      console.error("Login failed", error)
    }
  }

  return (
    <div style={styles.formContainer}>
      <form onSubmit={handleSubmit} style={styles.form}>
        <div style={styles.formGroup}>
          <div style={styles.inputWrapper}>
            <i className="user-icon" style={styles.icon}>
              👤
            </i>
            <input
              type="text"
              id="username"
              placeholder="Username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              style={styles.input}
            />
          </div>
        </div>
        <div style={styles.formGroup}>
          <div style={styles.inputWrapper}>
            <i className="lock-icon" style={styles.icon}>
              🔒
            </i>
            <input
              type="password"
              id="password"
              placeholder="Password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              style={styles.input}
            />
          </div>
        </div>
        {error && <p style={{ color: 'red' }}>{error}</p>}  {/* Показываем ошибку */}
        <button type="submit" style={styles.button}>
          Login
        </button>
        {/* Остальная часть формы */}
      </form>
    </div>
  )
}

const styles = {
  formContainer: {
    width: "100%",
    maxWidth: "400px",
    backgroundColor: "#fff",
    borderRadius: "16px",
    padding: "32px",
  },
  form: {
    display: "flex",
    flexDirection: "column",
    gap: "20px",
  },
  formGroup: {
    marginBottom: "15px",
  },
  inputWrapper: {
    position: "relative",
    display: "flex",
    alignItems: "center",
  },
  icon: {
    position: "absolute",
    left: "12px",
    color: "#666",
  },
  input: {
    width: "100%",
    padding: "12px 12px 12px 40px",
    borderRadius: "8px",
    border: "1px solid #ddd",
    fontSize: "16px",
    backgroundColor: "#f8f8f8",
  },
  button: {
    width: "100%",
    padding: "12px",
    backgroundColor: "#FF4B55",
    color: "#fff",
    border: "none",
    borderRadius: "8px",
    cursor: "pointer",
    fontSize: "16px",
    fontWeight: "500",
  },
   orLogin: {
      textAlign: "center",
      color: "#666",
      fontSize: "14px",
    },
    socialButtons: {
      display: "flex",
      justifyContent: "center",
      gap: "16px",
    },
    socialButton: {
      width: "40px",
      height: "40px",
      borderRadius: "8px",
      border: "1px solid #ddd",
      backgroundColor: "white",
      cursor: "pointer",
      display: "flex",
      alignItems: "center",
      justifyContent: "center",
    },
    register: {
      textAlign: "center",
      fontSize: "14px",
      color: "#666",
    },
    registerLink: {
      color: "#FF4B55",
      textDecoration: "none",
    }

}

export default LoginForm
