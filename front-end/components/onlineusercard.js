import { useRouter } from 'next/navigation';

export default function UserCard({ user }) {
    const router = useRouter();

    return (
        <div className="userCard gp12" onClick={() => router.push(`/chat/${user.id}`)}>
            {user.avatar ? (
                <img
                    className="avatar32"
                    src={`http://localhost:8080/${user.avatar}`}
                    alt={`${user.firstName}'s avatar`}
                />
            ) : (
                <div className="avatar32 fallback" />
            )}
            <div className="userInfo">
                <div className="userName">{user.firstName} {user.lastName}</div>
                <div className="followerCount">{user.status}</div>
            </div>
        </div>
    );
}

