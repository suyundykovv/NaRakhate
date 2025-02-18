"use client"

import { useState } from "react"
import { Link } from "react-router-dom"

const LoginForm = ({ onLogin }) => {
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [error, setError] = useState("") // State to hold error messages
  const [loading, setLoading] = useState(false) // State to manage loading state

  const handleSubmit = async (e) => {
    e.preventDefault()

    setLoading(true)
    setError("") // Clear previous errors

    const loginData = {
      email,
      password,
    }

    try {
      const response = await fetch("http://127.0.0.1:8080/log-in", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(loginData),
      })

      if (!response.ok) {
        throw new Error("Invalid credentials")
      }

      const data = await response.json()
      const token = data.token

      // If the login is successful, save the token and call onLogin
      localStorage.setItem("token", token) // You can store the token in localStorage for further authenticated requests.
      onLogin() // Call the onLogin prop function (this can be used to trigger a state update in the parent component)
    } catch (err) {
      setError("Login failed: " + err.message + loginData) // Display error message if login fails
    } finally {
      setLoading(false) // Stop loading when request is complete
    }
  }

  return (
    <div style={styles.formContainer}>
      <form onSubmit={handleSubmit} style={styles.form}>
        <div style={styles.formGroup}>
          <div style={styles.inputWrapper}>
            <i className="user-icon" style={styles.icon}>
              📧
            </i>
            <input
              type="email"
              id="email"
              placeholder="Email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
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
        {error && <div style={{ color: "red", textAlign: "center" }}>{error}</div>} {/* Show error message */}

        <button type="submit" style={styles.button} disabled={loading}>
          {loading ? "Logging in..." : "Login"}
        </button>

        <div style={styles.register}>
          Don't have an account yet?{" "}
          <Link to="/register" style={styles.registerLink}>
            Register
          </Link>
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

