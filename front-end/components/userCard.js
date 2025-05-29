import { useRouter } from 'next/navigation';
import { use, useState } from 'react';
import Avatar  from './avatar';

export default function UserCard({ user }) {
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
                //<div className="userAvatar fallback" />
                <Avatar user={user} />
            )}
            <div className="userInfo">
                <div className="userName">{user.firstName} {user.lastName}</div>
                <div className="followerCount">{user.followerCount} followers</div>
            </div>
            <button className={isFollowed ? "followedBtn" : "follow"} onClick={()=>{followReq(user.id)}}>
                {isFollowed ? "cancel follow" : "follow"}
            </button>
        </div>
    );
}
