"use client";

import { useState, useRef, useContext, useEffect } from 'react';
import { WebSocketContext } from '../context/wsContext';
import styles from './chatInput.module.css';

export default function MessageInput({ id, type }) {
  const [message, setMessage] = useState('');
  const { socket } = useContext(WebSocketContext);

  const [showEmojiModal, setShowEmojiModal] = useState(false);
  const modalRef = useRef(null);
  const emojis = ['😀', '😂', '😍', '😎', '😢', '🤔', '🔥', '🥳', '🙌', '💯'];


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

  const handleEmojiClick = (emoji) => {
    setMessage(prev => prev + emoji);
  };

  const handleClickOutside = (e) => {
    if (modalRef.current && !modalRef.current.contains(e.target)) {
      setShowEmojiModal(false);
    }
  };
  useEffect(() => {
    if (showEmojiModal) {
      document.addEventListener('mousedown', handleClickOutside);
    } else {
      document.removeEventListener('mousedown', handleClickOutside);
    }

    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [showEmojiModal]);


  return (
    <div className={styles.container }>
      <span
        onClick={() => setShowEmojiModal(prev => !prev)}
        className={styles.emojiToggle}
      >
        😊
      </span>
      {showEmojiModal && (
        <div className={styles.emojisSection} ref={modalRef}>
          {emojis.map((emoji) => (
            <span
              key={emoji}
              onClick={() => handleEmojiClick(emoji)}
              className={styles.emojiItem}
            >
              {emoji}
            </span>
          ))}
        </div>
      )}
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
