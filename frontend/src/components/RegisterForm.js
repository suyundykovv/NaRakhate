"use client"

import { useState } from "react"
import { useNavigate, Link } from "react-router-dom"

const RegisterForm = () => {
  const [username, setUsername] = useState("")
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [confirmPassword, setConfirmPassword] = useState("")
  const [error, setError] = useState("") // For error handling
  const [loading, setLoading] = useState(false) // For loading state management

  const navigate = useNavigate() // Initialize useNavigate hook

  const handleSubmit = async (e) => {
    e.preventDefault()

    if (password !== confirmPassword) {
      setError("Passwords do not match.")
      return
    }

    setLoading(true)
    setError("") // Reset previous errors

    const registrationData = {
      username,
      email,
      password,
      role: "user", // Default role for new users
    }

    try {
      const response = await fetch("http://localhost:8080/sign-up", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(registrationData),
      })

      if (!response.ok) {
        const errorText = await response.text() // Get error message from the response body
        throw new Error(`Registration failed: ${errorText}`)
      }

      const data = await response.json()
      alert("Registration successful! Please log in.")
      navigate("/login")
    } catch (err) {
      console.error("Error:", err) // Log detailed error
      setError("Registration failed: " + err.message) // Show error message in the UI
    } finally {
      setLoading(false) // Stop loading spinner
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
            <i className="email-icon" style={styles.icon}>
              ✉️
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
        <div style={styles.formGroup}>
          <div style={styles.inputWrapper}>
            <i className="lock-icon" style={styles.icon}>
              🔒
            </i>
            <input
              type="password"
              id="confirmPassword"
              placeholder="Confirm Password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              style={styles.input}
            />
          </div>
        </div>
        {error && <div style={{ color: "red", textAlign: "center" }}>{error}</div>} {/* Error Message */}
        <button type="submit" style={styles.button} disabled={loading}>
          {loading ? "Registering..." : "Register"}
        </button>
        <div style={styles.login}>
          Already have an account?{" "}
          <Link to="/login" style={styles.loginLink}>
            Login
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
  login: {
    textAlign: "center",
    fontSize: "14px",
    color: "#666",
  },
  loginLink: {
    color: "#FF4B55",
    textDecoration: "none",
  },
}

export default RegisterForm

