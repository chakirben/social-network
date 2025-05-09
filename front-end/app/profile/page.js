'use client'
import SideBar from "@/components/sidebar";
import "../profile/profile.css"
import { useEffect, useState } from "react";

export default function Profile() {
    const [profile, setProfile] = useState(null);

    useEffect(() => {
        async function fetchProfile() {
            try {
                const res = await fetch('http://localhost:8080/api/getUserData', { credentials: "include" });
                if (!res.ok) {
                    throw new Error(`HTTP error! Status: ${res.status}`);
                }
                const data = await res.json();
                setProfile(data);

            } catch (err) {
                console.error('Error fetching profile:', err);
            }
        }
        fetchProfile();
    }, []);

    return (
        <div className="profileContainer">
            <SideBar />
            {profile ? (
                <div className="profileCard">
                    <p><strong>Nickname:</strong> {profile.nickname}</p>
                    <p><strong>First Name:</strong> {profile.firstName}</p>
                    <p><strong>Last Name:</strong> {profile.lastName}</p>
                    <p><strong>Email:</strong> {profile.email}</p>
                    <p><strong>Age:</strong> {profile.age}</p>
                </div>
            ) : (
                <p>Loading...</p>
            )}


        </div>
    );
}
