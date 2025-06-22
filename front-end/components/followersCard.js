import { useRouter } from 'next/navigation';



export default function FollowersCard({ user }) {
    const router = useRouter();

    return (
        <div className="userCard" onClick={() => router.push(`/users/${user.ID}`)}>
            {user.avatar ? (
                <img
                    className="userAvatar"
                    src={`/${user.Avatar}`}
                    alt={`${user.Firstname}'s avatar`}
                />
            ) : (
                <div className="letterAvatar">
                    <span>{user.Firstname[0].toUpperCase()}{user.Lastname[0].toUpperCase()}</span>
                </div>
            )}
            <div className="userInfo">
                <div className="userName">{user.Firstname} {user.Lastname}</div>
            </div>
        </div>
    );
}
