'use client'
import SideBar from "@/components/sidebar";
import { useEffect, useState } from "react";
import "../../styles/global.css"
import "../profile/profile.css"


export default function Profile() {
    const [profile, setProfile] = useState(null)
    useEffect(() => {
        async function fetchProfile() {
            try {
                const res = await fetch('http://localhost:8080/api/getUserData', { credentials: "include" })
                if (!res.ok) {
                    throw new Error(`HTTP error! Status: ${res.status}`);
                }
                let data = await res.json()
                console.log(data);

                setProfile(data)
            } catch (err) {
                console.log('errrrror', err)
            }
        }
        fetchProfile()
    }, [])

    return (
        <div className="profileContainer">
            <SideBar />
            <div className="userProfile">
                <img className="coverture" src="http://localhost:8080/uploads/coverture.png"></img>
                <div className="userData">
                    <div className="imgAndFollow sb">
                        {profile && <img className="userAvatar" src={`http://localhost:8080/${profile.avatar}`} />}
                        <div className="follow">
                            <p><strong className="number">11K</strong> Followers</p>
                            <p><strong className="number">15K</strong> Following</p>
                        </div>
                    </div>
                    <div className="nameAndAbout">
                        {profile && <h4>{profile.firstName} {profile.lastName}</h4>}
                        {profile && <p>Hey everyone! I’ve been thinking about starting</p>}
                    </div>

                </div>

            <hr></hr>
            </div>


        </div>
    )
}

{/* {profile ? (
    <div className="profileCard">
        <p><strong>Nickname:</strong> {profile.nickname}</p>
        <p><strong>First Name:</strong> {profile.firstName}</p>
        <p><strong>Last Name:</strong> {profile.lastName}</p>
        <p><strong>Email:</strong> {profile.email}</p>
        <p><strong>Age:</strong> {profile.age}</p>
    </div>
) : (
    <p>Loading...</p>
)} */}