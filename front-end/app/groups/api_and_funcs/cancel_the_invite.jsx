"use client"

export default async function CancelTheInvite(userid , groupId) {
    try {
        const rep = await fetch(`/api/CancelInviteToGroups`
, {
            method: "POST",
            credentials: "include",
            headers: {
                "Content-Type": "application/json", 
            },
            body: JSON.stringify({
                userid : parseInt(userid,10),
                groupId: parseInt(groupId, 10)
            }),
        });
        if (!rep.ok) {
            throw new Error("Failed to send Request to join the group");
        }
      } catch (error) {
        console.error("Error sending the join request:", error.message || error);
    }
}