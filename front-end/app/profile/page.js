'use client'
import SideBar from "@/components/sidebar";
import { useEffect, useState } from "react";
import "../../styles/global.css"
import "../profile/profile.css"


export default function Profile() {
    const [profile, setProfile] = useState(null)
    const [showOptions, setShowOptions] = useState(false)


    const handle = () => {
        setShowOptions(!showOptions)
    }

    const handleSelection = async (choise) => {
        setShowOptions(false)
        try {
            let res = await fetch('http://localhost:8080/api/updatePrivacy', {
                method: 'POST',
                credentials: 'include',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ accountType: choise })
            })
            if (!res.ok) {
                throw new Error(`http error here : ${res.status}`)
            }
            setProfile(prev => ({ ...prev, accountType: choise }))
            console.log("choise", choise);
        } catch (error) {
            console.error("Failed to update privacy:", error)
        }
    }


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
                            <span onClick={handle} className="privacyOptions">
                                privacy {profile?.accountType === "public"
                                    ? <img className="icon" src="http://localhost:8080/uploads/unlock.svg" />
                                    : <img className="icon" src="http://localhost:8080/uploads/lock.svg" />}
                                <div className={`options ${showOptions ? "" : "hidden"}`}>
                                    <span className="button" onClick={() => handleSelection("private")}> <img className="icon" src="http://localhost:8080/uploads/lock.svg" /> Private</span>
                                    <span className="button" onClick={() => handleSelection("public")}> <img className="icon" src="http://localhost:8080/uploads/unlock.svg" /> Public</span>
                                </div>
                            </span>

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