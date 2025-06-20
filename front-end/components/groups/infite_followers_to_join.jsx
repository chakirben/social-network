'use client'
import { useEffect, useState } from "react"
import UserCard from "@/components/userCard"

export default function InviteFollowersToGroups({ onSkip, groupId }) {
    const [allFollowers, setAllFollowers] = useState([])

    useEffect(() => {
        const getFollowers = async () => {
            try {
                const rep = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/getFollowers`, {
                    method: "POST",
                    credentials: "include",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify({
                        groupId: parseInt(groupId, 10)
                    }),
                });

                if (!rep.ok) {
                    throw new Error("Failed to fetch the list of followers");
                }

                const repData = await rep.json();
                console.log(repData);
                setAllFollowers(repData);

            } catch (error) {
                console.error("An error occurred while fetching followers:", error);
            }
        };

        getFollowers();
    }, []);

    return (
        <div className="overlayy">
            <div className="AllFollowersAndSkip">
                <div className="AllFollowersToInvite">
                    {allFollowers && allFollowers.length > 0 ? (
                        allFollowers.map((user) => (
                            <UserCard key={user.id} user={user} invite="+invite" groupId={groupId} />
                        ))
                    ) : (
                        <div>No followers to display.</div>
                    )}
                </div>

                <div className="divcreatbtn">
                    <button className="skipbtn" onClick={onSkip}>
                        <img src="./../images/skip.svg" />
                        Skip
                    </button>
                </div>
            </div>
        </div>

    );
}
