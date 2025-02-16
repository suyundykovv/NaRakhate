import { useState } from "react"

const LoginForm = ({ onLogin }) => {
  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")

  const handleSubmit = (e) => {
    e.preventDefault()
    console.log("Login:", username)
    console.log("Password:", password)
    onLogin() // Вызов onLogin после успешной авторизации
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
        <div style={styles.rememberForgot}>
          <label style={styles.remember}>
            <input type="checkbox" /> Remember me
          </label>
          <a href="#" style={styles.forgot}>
            Forgot password?
          </a>
        </div>
        <button type="submit" style={styles.button}>
          Login
        </button>
        <div style={styles.orLogin}>or login with</div>
        <div style={styles.socialButtons}>
          <button style={styles.socialButton}>G</button>
          <button style={styles.socialButton}>f</button>
          <button style={styles.socialButton}>🍎</button>
        </div>
        <div style={styles.register}>
          Don't have an account yet?{" "}
          <a href="#C:\Users\Aday\Desktop\auth-page\src\components\RegisterForm.js" style={styles.registerLink}>
            Register
          </a>
        </div>
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
  rememberForgot: {
    display: "flex",
    justifyContent: "space-between",
    alignItems: "center",
    fontSize: "14px",
  },
  remember: {
    display: "flex",
    alignItems: "center",
    gap: "8px",
  },
  forgot: {
    color: "#666",
    textDecoration: "none",
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
  },
}

export default LoginForm