'use client'
import SideBar from "@/components/sidebar";
import { useEffect, useState } from "react";
import "../../styles/global.css"
import "../profile/profile.css"
import Post from "@/components/post";
import "../home/home.css"
import Header from "@/components/header";
import { useRouter } from "next/navigation";


export default function Profile() {
    const [profile, setProfile] = useState(null)
    const [profileData, setData] = useState([])
    const [showOptions, setShowOptions] = useState(false)
    const router = useRouter()

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
    const handleLogout = async () => {
        fetch('http://localhost:8080/api/logout', {
            method: 'POST',
            credentials: 'include',
        })
            .then((res) => {
                if (res.ok) {
                    console.log("logout");
                    router.push('/login')

                } else {
                    console.log("error");
                }
            })
            .catch((err) => {
                console.log("error", err);
            }
            )
    }

    useEffect(() => {
        async function fetchProfile() {
            try {
                const res = await fetch('http://localhost:8080/api/getUserData', { credentials: "include" })
                if (!res.ok) {
                    throw new Error(`HTTP error! Status: ${res.status}`);
                }
                let data = await res.json()
                // console.log(data);

                setProfile(data)
            } catch (err) {
                console.log('errrrror', err)
            }
        }
        fetchProfile()
    }, [])

    useEffect(() => {
        async function fetchPosts() {
            try {
                let res = await fetch('http://localhost:8080/api/GetCreatedPosts', {
                    method: 'POST',
                    credentials: 'include',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                })
                if (!res.ok) {
                    throw new Error(`http error here : ${res.status}`)
                }
                let posts = await res.json()
                setData(posts)
                console.log(posts);

            } catch (error) {
                console.error("errr", error)
            }
        }
        fetchPosts()
    }, [])



    return (
        <div className="profileContainer">
            <SideBar />
            <div className="classname df cl">
                <Header pageName={"profile"} ele={<button className="Secondary" onClick={handleLogout}>logout</button>} />
                <div className="userProfile">
                    <img className="coverture" src="http://localhost:8080/uploads/coverture.png"></img>
                    <div className="userdata gp12">
                        <div className="imgAndFollow sb">
                            {profile && <img className="userAvatar" src={`http://localhost:8080/${profile.avatar}`} />}
                            <div className="follow">
                                <p><strong className="number">{profile?.followers}</strong> Followers</p>
                                <p><strong className="number">{profile?.following}</strong> Following</p>
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
                            <div className="df gp6">
                                {profile && <img src="http://localhost:8080/uploads/dateOfBirth.svg" />}
                                {profile && <p>{new Date(profile.dateOfBirth).toLocaleDateString('fr-FR')}</p>}
                            </div>

                        </div>
                    </div>
                    <hr></hr>
                </div>

                {profileData && profileData.map((p, i) => (
                    <Post key={i} pst={p} />
                ))}
            </div>
        </div>
    )
}