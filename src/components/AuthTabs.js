import React, { useState } from 'react';
import LoginForm from './LoginForm';
import RegisterForm from './RegisterForm';

const AuthTabs = () => {
  const [activeTab, setActiveTab] = useState('login'); // Состояние для активной вкладки

  return (
    <div style={styles.container}>
      {/* Вкладки */}
      <div style={styles.tabs}>
        <button
          style={activeTab === 'login' ? styles.activeTab : styles.tab}
          onClick={() => setActiveTab('login')}
        >
          Вход
        </button>
        <button
          style={activeTab === 'register' ? styles.activeTab : styles.tab}
          onClick={() => setActiveTab('register')}
        >
          Регистрация
        </button>
      </div>

      {/* Форма входа */}
      {activeTab === 'login' && <LoginForm />}

      {/* Форма регистрации */}
      {activeTab === 'register' && <RegisterForm />}
    </div>
  );
};

const styles = {
  container: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    height: '100vh',
    backgroundColor: '#f0f0f0',
  },
  tabs: {
    display: 'flex',
    marginBottom: '20px',
  },
  tab: {
    padding: '10px 20px',
    border: 'none',
    backgroundColor: '#ddd',
    cursor: 'pointer',
    fontSize: '16px',
  },
  activeTab: {
    padding: '10px 20px',
    border: 'none',
    backgroundColor: '#007bff',
    color: '#fff',
    cursor: 'pointer',
    fontSize: '16px',
  },
};

export default AuthTabs;