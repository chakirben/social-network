'use client'
import { useState, useEffect } from "react";
import SideBar from "@/components/sidebar";
import UserCard from "@/components/userCard";

export default function Chat() {
  const [followedUsers, setFollowedUsers] = useState([]);

  useEffect(() => {
    const fetchFolowedUsers = async () => {
      try {
        const rs = await fetch("http://localhost:8080/api/getUnfollowedUsers", {
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
    <div>
     
      <p>Hello Chat</p>
      <ChatSection users={followedUsers } />
    </div>
  );
}
function  ChatSection({ users }) {
  return (
    
    <div className="usersList">
      {users.map(user => (
        <section key={user.id}  className="userCard">
          <UserCard user={user} Showchat = {true } />
        </section>
      ))}
    </div>
  );
}
