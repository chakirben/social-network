'use client';

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import SideBar from "@/components/sidebar";
import Post from "@/components/post";
import Header from "@/components/Header/header";
import { useUser } from "@/components/context/userContext";
import "../../../styles/global.css";
import "../../profile/profile.css";
import "../../home/home.css";
import Avatar from "@/components/avatar/avatar";

export default function ProfileClient({ session, searchParams }) {
  const [profileData, setProfileData] = useState(null);
  const [showFollowersModal, setShowFollowersModal] = useState(false);
  const [showFollowingModal, setShowFollowingModal] = useState(false);
  const [followersList, setFollowersList] = useState([]);
  const [followingList, setFollowingList] = useState([]);
  const [errorMessage, setErrorMessage] = useState(null);
  const [followersCount, setFollowersCount] = useState(null);
  const [followingCount, setFollowingCount] = useState(null);
  const [followBtnText, setFollowBtnText] = useState(null);

  const { id } = useParams() || searchParams;
  const router = useRouter();
  const { user } = useUser();

  useEffect(() => {
    if (profileData) {
      setFollowersCount(profileData.followers_count);
      setFollowingCount(profileData.followed_count);
      setFollowBtnText(profileData.follow_status); 
    }
  }, [profileData]);

  useEffect(() => {
    if (user?.id == id) {
      router.push('/profile');
    }
  }, [user, id]);

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const response = await fetch(`/api/profile`, {
          credentials: 'include',
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ session, id }),
        });

        if (!response.ok) {
          const errorText = await response.text();
          setErrorMessage(errorText || "Failed to load profile");
          setProfileData(null);
        } else {
          const data = await response.json();
          setProfileData(data);
          setErrorMessage(null);
        }
      } catch (error) {
        setErrorMessage("Unable to connect to server");
        setProfileData(null);
      }
    };
    if (id) fetchProfile();
  }, [id, session]);

  const fetchUserList = async (type) => {
    try {
      const response = await fetch(`/api/followersList?type=${type}&id=${id}`, {
        credentials: 'include',
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
      });

      const data = await response.json();
      if (type === 'followers') {
        setFollowersList(data);
        setFollowersCount(data?.length);
        setShowFollowersModal(true);
      } else {
        setFollowingList(data);
        setFollowingCount(data?.length);
        setShowFollowingModal(true);
      }
    } catch (error) {
      console.error("Fetch user list error:", error);
    }
  };
  const handleFollow = async (e) => {
    e.preventDefault();
    if (!followBtnText) return;

    try {
      const response = await fetch(`/api/follow?id=${id}&action=${followBtnText}`, {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          follower_session: session,
          followed_id: id,
        }),
      });

      if (!response.ok) {
        const errorText = await response.text();
        console.error('Error:', errorText);
        return;
      }

      const newStatus = await response.text();
      setFollowBtnText(newStatus);

      setFollowersCount((count) => {
        if (newStatus === "follow") {
          return count - 1 >= 0 ? count - 1 : 0;
        } else if (newStatus === "unfollow") {
          return count + 1;
        }
        return count;
      });

    } catch (error) {
      console.error("Follow request failed:", error);
    }
  };

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
    posts,
    profile_status,
    profile_type,
  } = profileData;

  return (
    <div className="profileContainer">
      <SideBar />
      <div className="classname df cl">
        <Header />
        <div className="userProfile">
          <img className="coverture" src="/images/coverture.png" />
          <div className="userdata gp12">
            <div className="imgAndFollow sb">
              <Avatar name={personal_data[0].Firstname} url={personal_data[0].Avatar} />

              <div className="follow">
                <p onClick={() => fetchUserList('followers')}>
                  <strong className="followers-number">{followersCount}</strong> Followers
                </p>
                <p onClick={() => fetchUserList('following')}>
                  <strong className="following-number">{followingCount}</strong> Following
                </p>
                {profile_status === "auther" && (
                  <button className="follow-button" onClick={handleFollow}>{followBtnText}</button>
                )}
              </div>
            </div>

            <h2>{personal_data[0].Nickname || `${personal_data[0].Firstname} ${personal_data[0].Lastname}`}</h2>
            <p>{personal_data[0].About || `${personal_data[0].Firstname}'s Profile`}</p>
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

        {showFollowersModal && (
          <div className="modal-backdrop" onClick={() => setShowFollowersModal(false)}>
            <div onClick={e => e.stopPropagation()}>
              <h2>Followers</h2>
              {followersList?.length > 0 ? (
                followersList.map(user => (
                  <div key={user.id || user.ID} className="df gp6 center" style={{ gap: '10px', marginBottom: '8px' }}>
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
            <div onClick={e => e.stopPropagation()}>
              <h2>Following</h2>
              {followingList?.length > 0 ? (
                followingList.map(user => (
                  <div key={user.id || user.ID} className="df gp6 center" style={{ gap: '10px', marginBottom: '8px' }}>
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
      </div>
    </div>
  );
}