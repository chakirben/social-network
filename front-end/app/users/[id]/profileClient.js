'use client';
import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import FollowButton from './FollowButton';
import SideBar from "@/components/sidebar";
import Post from "@/components/post";
import Header from "@/components/Header/header";
import { useUser } from "@/components/context/userContext";
import "../../../styles/global.css";
import "../../profile/profile.css";
import "../../home/home.css";
import FollowersCard from "@/components/followersCard";
import Avatar from "@/components/avatar/avatar";

export default function ProfileClient({ session, searchParams }) {
    const [profileId, setProfileId] = useState(null);
    const [profileData, setProfileData] = useState(null);
    const [showFollowModal, setShowFollowModal] = useState(false);
    const [errorMessage, setErrorMessage] = useState(null);
    const { id } = useParams() || searchParams;

    const route = useRouter()
    const { user } = useUser();

    useEffect(() => {
        if (user?.id == id) {
            route.push('/profile');
        }
    }, [user, id]);


    useEffect(() => {
        setProfileId(id);
    }, [id]);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await fetch('http://localhost:8080/api/profile', {
                    credentials: 'include',
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ session: session, id: id }),
                });

                if (!response.ok) {
                    const errorText = await response.text();
                    setErrorMessage(errorText || "Failed to load profile");
                    setProfileData(null);
                    console.error("Server error:", errorText);
                } else {
                    const data = await response.json();
                    setProfileData(data);
                    setErrorMessage(null); // clear error
                    console.log("Success dasetShowFollowModalta:", data);
                }
            } catch (error) {
                console.error("Network error:", error);
                setErrorMessage("Unable to connect to server");
                setProfileData(null);
            }
        };

        if (id) fetchData();
    }, [id]);

    // Show error or loading state
    if (!profileData) {
        return (
            <div className="profileContainer">
                <SideBar />
                <div className="classname df cl">
                    <Header />
                    <div className="userProfile" style={{ padding: "2rem", textAlign: "center" }}>
                        {errorMessage ? (
                            <div style={{ color: "gray", fontSize: "1.1rem" }}>⚠️ {errorMessage}</div>
                        ) : (
                            <div>Loading profile...</div>
                        )}
                    </div>
                </div>
            </div>
        );
    }

    const {
        personal_data,
        followers_count,
        followed_count,
        followers_data,
        followeds_data,
        posts,
        follow_status,
        profile_status,
        profile_type
    } = profileData;

    return (
        <div className="profileContainer">
            <SideBar />
            <div className="classname df cl">
                <Header />
                <div className="userProfile">
                    <img className="coverture" src={"http://localhost:8080/uploads/coverture.png"} alt="Coverture" />
                    <div className="userdata gp12">
                        <div className="imgAndFollow sb">

                            <Avatar name={personal_data[0].Firstname} url={personal_data[0].Avatar} />


                            <div className="follow">
                                <p onClick={() => setShowFollowModal(true)}><strong className="followers-number">{followers_count}</strong> Followers</p>
                                <p onClick={() => setShowFollowModal(true)}><strong className="following-number">{followed_count}</strong> Following</p>
                                {profile_status === "auther" && (
                                    <FollowButton follow_status={follow_status} session={session} id={id} />
                                )}
                            </div>
                        </div>

                        <h2>{personal_data[0].Nickname || personal_data[0].Firstname + " " + personal_data[0].Lastname}</h2>
                        <p>{personal_data[0].About || personal_data[0].Firstname + "'s Profile"}</p>
                    </div>
                </div>

                {profile_type === "private" && (
                    <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', gap: '8px', marginTop: '16px', color: '#555' }}>
                        <span>🔒</span>
                        <span>This is a private account</span>
                    </div>
                )}

                {posts && posts.map((p, i) => (
                    <Post key={i} pst={p} />
                ))}
            </div>
            {showFollowModal && (
                <div className="modal-backdrop" onClick={() => setShowFollowModal(false)} style={{
                    position: 'fixed',
                    top: 0,
                    left: 0,
                    width: '100vw',
                    height: '100vh',
                    backgroundColor: 'rgba(0, 0, 0, 0.5)',
                    display: 'flex',
                    justifyContent: 'center',
                    alignItems: 'center',
                    zIndex: 1000
                }}>
                    <div
                        onClick={(e) => e.stopPropagation()}
                        style={{
                            backgroundColor: 'black',
                            padding: '20px',
                            borderRadius: '8px',
                            maxHeight: '80vh',
                            overflowY: 'auto',
                            width: '650px',
                        }}
                    >
                        <div className="followers_modal">
                            <h2>Followers users</h2>
                            {followers_data && followers_data.length > 0 ? (
                                followers_data.map((user) => (
                                    <FollowersCard key={user.ID} user={user} />
                                ))
                            ) : (
                                <div>There is no followers</div>
                            )}
                        </div>
                        <div className="followeds_modal">
                            <h2>Following users</h2>
                            {followeds_data && followeds_data.length > 0 ? (
                                followeds_data.map((user) => (
                                    <FollowersCard key={user.ID} user={user} />
                                ))
                            ) : (
                                <div>There is no following</div>
                            )}
                        </div>
                        <button onClick={() => setShowFollowModal(false)} style={{ marginTop: '10px' }}>
                            Close
                        </button>
                    </div>
                </div>
            )}
        </div>
    );
}
