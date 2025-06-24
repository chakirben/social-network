'use client'
import SideBar from "@/components/sidebar";
import { useContext, useEffect, useState } from "react";
import "../../styles/global.css"
import "../profile/profile.css"
import Post from "@/components/post";
import "../home/home.css"
import Header from "@/components/Header/header";
import { useRouter } from "next/navigation";
import Divider from "@/components/divider";
import { WebSocketContext } from "@/components/context/wsContext";
import Avatar from "@/components/avatar/avatar";


export default function Profile() {
    const [profile, setProfile] = useState(null)
    const [profileData, setData] = useState([])
    const [showOptions, setShowOptions] = useState(false)
    const router = useRouter()
    const socket = useContext(WebSocketContext)
    const [showFollowersModal, setShowFollowersModal] = useState(false);
    const [showFollowingModal, setShowFollowingModal] = useState(false);
    const [userList, setUserList] = useState([]);
    console.log(socket);
    const fetchUserList = async (listType) => {
        try {
            const response = await fetch(`/api/followersList?&type=${listType}`, {
                credentials: 'include',
            });
            if (!response.ok) {
                console.error("Error:", await response.text());
                return;
            }
            const data = await response.json();
            setUserList(data);
            if (listType === "followers") {
                setShowFollowersModal(true);
            } else {
                setShowFollowingModal(true);
            }
        } catch (error) {
            console.error(error);
        }
    };
    const handle = () => {
        setShowOptions(!showOptions)
    }

    const handleSelection = async (choise) => {
        setShowOptions(false)
        try {
            let res = await fetch(`/api/updatePrivacy`, {
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
        fetch(`/api/logout`, {
            method: 'POST',
            credentials: 'include',
        })
            .then((res) => {
                if (res.ok) {
                    console.log("logout")
                    localStorage.removeItem("user")
                    router.push('/login')
                    socket.socket.close()

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
                const res = await fetch(`/api/getUserData`, { credentials: "include" })
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
                let res = await fetch(`/api/GetCreatedPosts`, {
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
                    <img className="coverture" src="images/coverture.png"></img>
                    <div className="userdata gp12">
                        <div className="imgAndFollow sb">
                            <Avatar url={profile?.avatar} size={"big"} name={profile?.firstName} />
                            <div className="follow">
                                <p onClick={() => fetchUserList("followers")}><strong className="number">{profile?.followers}</strong> Followers</p>
                                <p onClick={() => fetchUserList("following")}><strong className="number">{profile?.following}</strong> Following</p>
                                <span onClick={handle} className="privacyOptions">
                                    privacy {profile?.accountType === "public"
                                        ? <img className="icon" src="images/unlock.svg" />
                                        : <img className="icon" src="images/lock.svg" />}
                                    <div className={`options ${showOptions ? "" : "hidden"}`}>
                                        <span className="button" onClick={() => handleSelection("private")}> <img className="icon" src="images/lock.svg" /> Private</span>
                                        <span className="button" onClick={() => handleSelection("public")}> <img className="icon" src="images/unlock.svg" /> Public</span>
                                    </div>
                                </span>

                            </div>
                        </div>
                        <div className="nameAndAbout">
                            {profile && <h4>{profile.firstName} {profile.lastName}</h4>}
                            {profile && <p><strong>Username: </strong> {profile.nickname}</p>}
                            {profile && <p>{profile.about}</p>}
                            {profile && <p>{profile.email}</p>}
                            <div className="df gp6">
                                {profile && <img src="images/dateOfBirth.svg" />}
                                {profile && <p>{new Date(profile.dateOfBirth).toLocaleDateString('fr-FR')}</p>}
                            </div>

                        </div>
                    </div>
                    <Divider />
                </div>
                {showFollowersModal && (
                    <div className="modal-backdrop" onClick={() => setShowFollowersModal(false)}>
                        <div onClick={(e) => e.stopPropagation()} >
                            <h2>Followers</h2>
                            {userList?.length > 0 ? (
                                userList.map((user) => (
                                    <div key={user.id} className="df gp6 center">
                                        <Avatar url={user.avatar} name={user.firstName} />
                                        <span>{user.firstName} {user.lastName}</span>
                                    </div>
                                ))
                            ) : (
                                <div>No followers</div>
                            )}
                            <button onClick={() => setShowFollowersModal(false)}>Close</button>
                        </div>
                    </div>
                )}
                {showFollowingModal && (
                    <div className="modal-backdrop" onClick={() => setShowFollowingModal(false)}>
                        <div onClick={(e) => e.stopPropagation()} >
                            <h2>Following</h2>
                            {userList?.length > 0 ? (
                                userList.map((user) => (
                                    <div key={user.id}  className="df gp6 center">
                                        <Avatar url={user.avatar} name={user.firstName} />
                                        <span>{user.firstName} {user.lastName}</span>
                                    </div>
                                ))
                            ) : (
                                <div>Not following anyone</div>
                            )}
                            <button onClick={() => setShowFollowingModal(false)}>Close</button>
                        </div>
                    </div>
                )}

                {profileData && profileData.map((p, i) => (
                    <Post key={i} pst={p} />
                ))}
            </div>
        </div>
    )
}