import { useEffect, useState, useRef } from 'react';
import { useParams } from 'react-router-dom';

const ChatWindow = ({ baseUrl, token, userId }) => {
  const { otherId } = useParams();
  const [messages, setMessages] = useState([]);
  const [text, setText] = useState('');
  const ws = useRef(null);
  const messagesEndRef = useRef(null);

  useEffect(() => {
    fetchMessages();
    connectWebSocket();

    return () => {
      if (ws.current) ws.current.close();
    };
  }, [otherId]);

  const fetchMessages = async () => {
    try {
      const res = await fetch(`${baseUrl}/messages?with=${otherId}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      const data = await res.json();
      setMessages(data);
      scrollToBottom();
    } catch (err) {
      console.error('Error fetching messages');
    }
  };

  const connectWebSocket = () => {
    ws.current = new WebSocket(`${(location.protocol === 'https:' ? 'wss' : 'ws')}://${location.host}/api/ws?user_id=${userId}`);

    ws.current.onopen = () => console.log('WebSocket connected');

    ws.current.onmessage = (event) => {
      const msg = JSON.parse(event.data);
      setMessages((prev) => [...prev, msg]);
      scrollToBottom();
    };

    ws.current.onclose = () => console.log('WebSocket closed');
  };

  const sendMessage = async () => {
    if (!text) return;
    const payload = { to_user_id: parseInt(otherId), text };
    try {
      await fetch(`${baseUrl}/messages`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(payload),
      });
      setText('');
    } catch (err) {
      console.error('Error sending message');
    }
  };

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  return (
    <div className="chat-window">
      <div className="messages">
        {messages.map((msg) => (
          <div key={msg.id} className={`message ${msg.sender_id === userId ? 'self' : ''} new`}>
            {msg.text}
          </div>
        ))}
        <div ref={messagesEndRef} />
      </div>
      <div className="input-container">
        <input
          type="text"
          value={text}
          onChange={(e) => setText(e.target.value)}
          onKeyDown={(e) => e.key === 'Enter' && sendMessage()}
          placeholder="Type a message..."
        />
        <button onClick={sendMessage}>Send</button>
      </div>
    </div>
  );
};

export default ChatWindow;