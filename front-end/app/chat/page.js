'use client'
import { useState, useEffect } from "react";
import SideBar from "@/components/sidebar";
import UserCard, { UserCardChat } from "@/components/userCard";
import Header from "@/components/header";

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
    <div className="home">
      <SideBar />
      
    <div className="homeP">
      <Header pageName={'chat'}></Header>
      <ChatSection users={followedUsers } />
    </div>
  </div>
  );
}
function  ChatSection({ users }) {
  return (
    <div className="usersList">
      {users.map(user => (
        <section key={user.id}  className="userCard">
          <UserCardChat user={user} Showchat = {true } />
        </section>
      ))}
    </div>
  );
}
