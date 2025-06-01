"use client";

import { useState, useContext } from 'react';
import { WebSocketContext } from './context/wsContext';
export default function MessageInput({ conversationId }) {
    const [message, setMessage] = useState('');
    const { ws } = useContext(WebSocketContext);

    const sendMessage = () => {
        if (message.trim() && ws && conversationId) {
            // Send message through WebSocket
            ws.send(JSON.stringify({
                type: 'message',
                content: message.trim(),
                conversationId: conversationId
            }));
            
            // Clear input
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
        <div className="message-input-container">
            <input
                type="text"
                value={message}
                onChange={(e) => setMessage(e.target.value)}
                onKeyPress={handleKeyPress}
                placeholder="Type your message"
                className="message-input"
            />
            <button 
                onClick={sendMessage}
                className="send-button"
            >
                <svg 
                    width="20" 
                    height="20" 
                    viewBox="0 0 24 24" 
                    fill="none" 
                    xmlns="http://www.w3.org/2000/svg"
                >
                    <path 
                        d="M2 21L23 12L2 3V10L17 12L2 14V21Z" 
                        fill="white"
                    />
                </svg>
            </button>
        </div>
    );
}