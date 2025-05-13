import { useRouter } from 'next/navigation';
import { use, useState } from 'react';

export default function UserCard({ user , Showchat = false}) {
    const router = useRouter();
    const [isFollowed, setIsfollowed] = useState(user.requested)
    const followReq = async (id) => {
        try {
            const response = await fetch(`http://localhost:8080/api/follow?id=${id}`, { credentials: "include" });
            if (response.ok) {
                if (isFollowed) {
                    setIsfollowed(false)
                } else {
                    setIsfollowed(true)
                }
            }
        } catch (error) {
            console.error("Error fetching posts:", error);
        }
    }
    return (
        <div className="userCard" onClick={() => router.push(`/users/${user.id}`)}>
            {user.avatar ? (
                <img
                    className="userAvatar"
                    src={`http://localhost:8080/${user.avatar}`}
                    alt={`${user.firstName}'s avatar`}
                />
            ) : (
                <div className="userAvatar fallback" />
            )}
            <div className="userInfo">
                <div className="userName">{user.firstName} {user.lastName}</div>
                {!Showchat ? <div className="userEmail">{user.email}</div> : null}
            </div>
            {!Showchat ? (
        <button className="followButton" onClick={() => followReq(user.id)}>
          {isFollowed ? "Unfollow" : "Follow"}
        </button>
      ) : (
        <p>Send</p>
      )}
        </div>
    );
}
