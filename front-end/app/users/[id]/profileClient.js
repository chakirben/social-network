'use client'
import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import FollowButton from './FollowButton';
import SideBar from "@/components/sidebar";
import Post from "@/components/post";
import Header from "@/components/Header/header";
import "../../../styles/global.css";
import "../../profile/profile.css";
import "../../home/home.css";
import Divider from "@/components/divider";

export default function ProfileClient({ session, searchParams }) {
    const [profileId, setProfileId] = useState(null);
    const [profileData, setProfileData] = useState(null);
    const [showFollowModal, setShowFollowModal] = useState(false);

    const { id } = useParams() || searchParams;

    // Set profile ID
    useEffect(() => {
        setProfileId(id);
    }, [id]);

    // Fetch profile data
    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await fetch('http://localhost:8080/api/profile', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ session: session, id: id }),
                });

                if (!response.ok) {
                    const errorText = await response.text();
                    console.error('Error fetching data:', errorText);
                } else {
                    const data = await response.json();
                    setProfileData(data); // <-- Save it to state
                    console.log("success data :", data);
                }
            } catch (error) {
                console.error('Network error:', error);
            }
        };

        if (id) fetchData();
    }, [id]);

    // Show loading state if data not loaded yet
    if (!profileData) {
        return <div>Loading profile...</div>;
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
                            <img className="userAvatar" src={ personal_data[0].Avatar || "http://localhost:8080/default-avatar.png"} alt="Avatar" />
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
                            backgroundColor: 'white',
                            padding: '20px',
                            borderRadius: '8px',
                            maxHeight: '80vh',
                            overflowY: 'auto',
                            width: '400px'
                        }}
                    >
                        <div className="followers_modal">
                            <h2 style={{ color: 'black', }}>Followers users</h2>
                            <ul>
                                {followers_data && followers_data.map((user) => (
                                    <li key={user.ID} style={{ cursor: 'pointer' }}>
                                        <a href={`/users/${user.ID}`}>
                                            {user.Firstname} {user.Lastname}
                                        </a>
                                    </li>
                                ))}
                                {!followers_data && <div style={{ color: 'black', }}>There is no followers</div>}
                            </ul>
                        </div>
                        <div className="followeds_modal">
                            <h2 style={{ color: 'black', }}>Following users</h2>
                            <ul>
                                {followeds_data && followeds_data.map((user) => (
                                    <li key={user.ID} style={{ cursor: 'pointer' }}>
                                        <a href={`/users/${user.ID}`}>
                                            {user.Firstname} {user.Lastname}
                                        </a>
                                    </li>
                                ))}
                                {!followeds_data && <div style={{ color: 'black', }}>There is no following</div>}
                            </ul>
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
