"use client"

export default async function FetchJoinToGroup(groupId) {
    try {
        const rep = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/RequestToJoinGroups`
, {
            method: "POST",
            credentials: "include",
            headers: {
                "Content-Type": "application/json", 
            },
            body: JSON.stringify({
                groupId: groupId, 
            }),
        });

        if (rep.status !== 202) {
            throw new Error("Failed to send Request to join the group");
        }
       
    } catch (error) {
        console.error("Error sending the join request:", error.message || error);
    }
}
