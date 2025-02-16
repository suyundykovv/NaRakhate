// api.js
const API_BASE_URL = 'http://127.0.0.1:8080'; // Замените на ваш базовый URL бэкенда

export const fetchBets = async () => {
  try {
    const response = await fetch(`${API_BASE_URL}/getAllBets`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });
    if (!response.ok) {
      throw new Error('Failed to fetch bets');
    }
    return await response.json();
  } catch (error) {
    console.error('Error fetching bets:', error);
    throw error;
  }
};

export const createBet = async (betData) => {
  try {
    const response = await fetch(`${API_BASE_URL}/createBet`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(betData),
    });
    if (!response.ok) {
      throw new Error('Failed to create bet');
    }
    return await response.json();
  } catch (error) {
    console.error('Error creating bet:', error);
    throw error;
  }
};

export const fetchMatchesByLeague = async (leagueId) => {
  try {
    const response = await fetch(`${API_BASE_URL}/fetchLeagueMatches?league=${leagueId}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });
    if (!response.ok) {
      throw new Error('Failed to fetch matches');
    }
    return await response.json();
  } catch (error) {
    console.error('Error fetching matches:', error);
    throw error;
  }
};

