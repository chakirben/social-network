"use client"

export default async function InviteTheFollowers(userid , groupId) {

    try {
        const rep = await fetch(`/api/InfiteTheFollowers`, {
            method: "POST",
            credentials: "include",
            headers: {
                "Content-Type": "application/json", 
            },
            body: JSON.stringify({
                userid: userid,
                groupId : parseInt(groupId, 10)
            }),
        });

        if (!rep.ok) {
            throw new Error("Failed to send Invite to join the group");
        }
       
    } catch (error) {
        console.error("Error sending the join request:", error.message || error);
    }
}