"use client";

import { useState, useContext } from 'react';
import { WebSocketContext } from '../context/wsContext';
import styles from './chatInput.module.css';

export default function MessageInput({ id, type }) {
  const [message, setMessage] = useState('');
  const { socket } = useContext(WebSocketContext);

  const sendMessage = () => {
    if (message.trim() && socket && id && type) {
      const parsedId = parseInt(id);
      const msgType = type === "group" ? "groupmsg" : "message";

      const payload = {
        type: msgType,
        content: message.trim(),
        ...(type === "group"
          ? { groupId: parsedId }
          : { receiverId: parsedId }) // for private messages
      };

      socket.send(JSON.stringify(payload));
      setMessage('');
    }
  };

  const handleKeyPress = (e) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      sendMessage();
    }
  };

  return (
    <div className={styles.container}>
      <input
        type="text"
        value={message}
        onChange={(e) => setMessage(e.target.value)}
        onKeyPress={handleKeyPress}
        placeholder="What's happening?"
        className={styles.input}
      />
      <button 
        onClick={sendMessage}
        className={styles.button}
      >
        Send
      </button>
    </div>
  );
}
