import { useRouter } from 'next/navigation';
import { useState } from 'react';
import Avatar from './avatar/avatar';
import InviteTheFollowers from "@/app/groups/api_and_funcs/infite_the_followers"
import CancelTheInvite from "@/app/groups/api_and_funcs/cancel_the_invite"

export default function UserCard({ user, invite, groupId }) {
    const router = useRouter();

    const [hasRequested, setHasRequested] = useState(user.hasRequested);
    const [isFollowed, setIsFollowed] = useState(user.isFollowed);

    const [clickToInvite ,  setclickToInvite] = useState(user.status)

    const getInitialBtnText = () => {
        if (hasRequested) return "cancel_request";
        if (isFollowed) return "unfollow";
        return "follow";
    };

    const [btnText, setBtnText] = useState(getInitialBtnText());

    const followReq = async (id, action) => {
        try {
            const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/follow?id=${id}&action=${action}`, {
                credentials: "include",
            });

            if (response.ok) {
                const newText = await response.text();

                if (newText === "follow") {
                    setIsFollowed(false);
                    setHasRequested(false);
                } else if (newText === "cancel_request") {
                    setHasRequested(true);
                    setIsFollowed(false);
                } else if (newText === "unfollow") {
                    setIsFollowed(true);
                    setHasRequested(false);
                }

                setBtnText(newText);
            }
        } catch (error) {
            console.error("Error sending follow request:", error);
        }
    };

    const handlerToInvite = () => {
        setclickToInvite("Cancel-Invite")
        InviteTheFollowers(user.id, groupId)
    }

    const handlerCancelInvite = () => {
        setclickToInvite("+invite")
        CancelTheInvite(user.id, groupId)
    }

    return (
        <>
            {invite === "+invite" ? (
                <div className="userCard">
                    <Avatar url={user.avatar} name={user.firstName} />
                    <div className="userInfo">
                        <div className="userName">{user.firstName} {user.lastName}</div>
                        <div className="followerCount">{user.followerCount} followers</div>
                    </div>
                    {clickToInvite == "+invite" ? (
                        <button onClick={handlerToInvite} >{clickToInvite}</button>

                    ) : (
                        <button onClick={handlerCancelInvite} >{clickToInvite}</button>
                    )}
                </div>
            ) : (
                <div className="userCard" onClick={() => router.push(`/users/${user.id}`)}>
                    <Avatar url={user.avatar} name={user.firstName} />
                    <div className="userInfo">
                        <div className="userName">{user.firstName} {user.lastName}</div>
                        <div className="followerCount">{user.followerCount} followers</div>
                    </div>

                    <>
                        <button
                            className={
                                btnText === "unfollow" || btnText === "cancel_request"
                                    ? "followedBtn"
                                    : "follow"
                            }
                            onClick={(e) => {
                                e.stopPropagation();
                                followReq(user.id, btnText);
                            }}
                        >
                            {btnText}
                        </button>
                    </>
                </div>
            )}


        </>

    );
}
