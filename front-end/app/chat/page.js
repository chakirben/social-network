'use client'
import { useState, useEffect } from "react";
import SideBar from "@/components/sidebar";
import UserCard, { UserCardChat } from "@/components/userCard";
import Header from "@/components/header";
import { useRouter } from 'next/navigation';

export default function Chat() {
  let  [followedUsers, setFollowedUsers] = useState([]);

  useEffect(() => {
    const fetchFolowedUsers = async () => {
      try {
        const rs = await fetch("http://localhost:8080/api/getFollowedUsers", {
          credentials: "include"
        });
        const data = await rs.json();
        setFollowedUsers(data); 
      } catch (error) {
        console.error("Error fetching posts:", error);
      }
    };

    fetchFolowedUsers();
  }, []);

  return (
  
    <div className="home">
      <SideBar />
      
    <div className="homeP">
      <Header pageName={'chat'}></Header>
      {followedUsers.length === 0 ? (
        <section className="no-friends-section">
          <div className="no-friends-content">
            <img src="/images/best-friend.png" alt="No friends" className="w-24 h-24 mb-4 opacity-80" />
            <h2 className="text-xl font-semibold text-gray-700">No friends</h2>
            <p className="text-gray-500 mt-2">Please add some friends to start a conversation</p>
          </div>
</section>
      ) : (
        <ChatSection users={followedUsers } />
      )}
    </div>
  </div>
  );
}

function ChatSection({ users }) {
  const [Chatiseractive, setChatActive] = useState(null);
  const [messages, setMessages] = useState([]);

  const handleSend = (text) => {
    const newMsg = {
      sender: "You",
      text,
      time: new Date().toLocaleTimeString(),
      fromUser: true
    };
    setMessages([...messages, newMsg]);
  };

  if (Chatiseractive) {
    return (
      <div className="usersList">
      <div className="chat-container">

        <h2 className="label">Chat with {Chatiseractive.firstName}</h2>
                <button className ="btback" onClick={() => setChatActive(null)} >← Back to friends</button>
        <ChatBox messages={messages} onSend={handleSend} />
      </div>
      </div>
    );
  }

  return (
    <div className="usersList">
      {users.map(user => (
        <div key={user.id} className="userCard">
          <UserCardChat user={user}  onClick={() => setChatActive(user)} />
        </div>
      ))}
    </div>
  );
}

import { useRef } from "react";

export function ChatBox({ messages = [], onSend }) {
  const [input, setInput] = useState("");
  const messagesEndRef = useRef(null);

  useEffect(() => {
    if (messagesEndRef.current) {
      messagesEndRef.current.scrollIntoView({ behavior: "smooth" });
    }
  }, [messages]);

  const handleSend = () => {
    if (input.trim() === "") return;
    onSend(input.trim());
    setInput("");
  };

  const handleKeyPress = (e) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  return (
    <div>
      <div className="chat-box" >
        {messages.length === 0 ? (
          <p>No messages yet</p>
        ) : (
          messages.map((msg, i) => (
            <div className="chat-message" key={i}>
              <strong className="message-sender">{msg.sender}</strong>
              {msg.text}
            <div className="time">{msg.time}</div>
            <div className="chat-me" > {msg.fromUser}</div>
            </div>
            
          ))
        )}
        <div ref={messagesEndRef} />
      </div>
      <div className="chat-input-container">
        <textarea
              className="message-input" 
              value={input}
              onChange={(e) => setInput(e.target.value)}
              onKeyDown={handleKeyPress}
              placeholder="Type your message..."
              rows={2}
            />
      </div>
      
      
      <button onClick={handleSend}>Send</button>
      
    </div>
  );
}

