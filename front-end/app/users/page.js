'use client'
import Header from "@/components/header";
import SideBar from "@/components/sidebar";
import "../home/home.css"
import { use, useState, useEffect } from "react";
import UserCard from "@/components/userCard";
export default function Users() {
    const [users, setUsers] = useState([]);
    useEffect(() => {
        const fetchUsers = async () => {
            try {
                const response = await fetch("http://localhost:8080/api/getUnfollowedUsers", { credentials: "include" });
                const data = await response.json();
                setUsers(data);
            } catch (error) {
                console.error("Error fetching posts:", error);
            }
        };
        fetchUsers();
    }, []);
    return (
        <div className="home">
            <SideBar />
            <div className="homeP">
                <Header pageName={'users'}></Header>
                <div className="usersList">
                    {users?.map(user => (
                        <UserCard key={user.id} user={user}  />
                    ))}
                </div>
            </div>

        </div>
    )
}