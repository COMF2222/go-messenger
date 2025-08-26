import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';

const ChatSidebar = ({ baseUrl, token }) => {
  const [interlocutors, setInterlocutors] = useState([]);
  const [onlineStatus, setOnlineStatus] = useState({});

  useEffect(() => {
    fetchInterlocutors();
  }, []);

  const fetchInterlocutors = async () => {
    try {
      const res = await fetch(`${baseUrl}/interlocutors`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      const data = await res.json();
      setInterlocutors(data);
      data.forEach(async (user) => {
        const statusRes = await fetch(`${baseUrl}/online/${user.id}`);
        const statusData = await statusRes.json();
        setOnlineStatus((prev) => ({ ...prev, [user.id]: statusData.online }));
      });
    } catch (err) {
      console.error('Error fetching interlocutors');
    }
  };

  return (
    <div className="sidebar">
      <h3>Chats</h3>
      {interlocutors.map((user) => (
        <Link to={`/chat/${user.id}`} key={user.id}>
          <div className="user-item">
            {user.username} <span className={onlineStatus[user.id] ? 'online' : 'offline'}>{onlineStatus[user.id] ? 'Online' : 'Offline'}</span>
          </div>
        </Link>
      ))}
    </div>
  );
};

export default ChatSidebar;