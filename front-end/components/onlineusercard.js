import { useRouter } from 'next/navigation';
import Avatar from './avatar/avatar';

export default function UserCard({ user }) {
    const router = useRouter();

    return (
        <div className="userCard gp12" onClick={() => router.push(`/chat/user${user.id + "_" + user.firstName + "_" + user.lastName}`)}>
            <Avatar name={user.firstName} url={user.avatar} />

            <div className="userInfo">
                <div className="userName">{user.firstName} {user.lastName}</div>
                <div className="followerCount">{user.status}</div>
            </div>
        </div>
    );
}

