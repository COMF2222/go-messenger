import { useEffect, useState } from 'react';
import { Routes, Route, Navigate, useNavigate } from 'react-router-dom';
import Login from './components/Login.jsx';
import Register from './components/Register.jsx';
import ChatSidebar from './components/ChatSidebar.jsx';
import ChatWindow from './components/ChatWindow.jsx';

const BASE_URL = '/api';

function App() {
  const [token, setToken] = useState(localStorage.getItem('token'));
  const [userId, setUserId] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    if (token) {
      const payload = JSON.parse(atob(token.split('.')[1]));
      setUserId(payload.user_id);

      fetch(`${BASE_URL}/interlocutors`, {
        headers: { Authorization: `Bearer ${token}` },
      })
        .then(res => {
          if (res.status === 401) {
            localStorage.removeItem('token');
            setToken(null);
            navigate('/login');
            return null;
          }
          return res.json();
        })
        .then(data => {
          if (!data) return;
          if (Array.isArray(data) && data.length > 0) {
            navigate(`/chat/${data[0].id}`);
          } else {
            navigate('/');
          }
        })
        .catch(() => navigate('/'));
    }
  }, [token, navigate]);


  if (!token) {
    return (
      <Routes>
        <Route path="/login" element={<Login setToken={setToken} />} />
        <Route path="/register" element={<Register setToken={setToken} />} />
        <Route path="*" element={<Navigate to="/login" />} />
      </Routes>
    );
  }

  return (
    <div className="glass-container">
      <ChatSidebar baseUrl={BASE_URL} token={token} />
      <Routes>
        <Route path="/" element={<div>Select a chat</div>} />
        <Route
          path="/chat/:otherId"
          element={<ChatWindow baseUrl={BASE_URL} token={token} userId={userId} />}
        />
        <Route path="*" element={<div>Select a chat</div>} />
      </Routes>
    </div>
  );
}

export default App;
