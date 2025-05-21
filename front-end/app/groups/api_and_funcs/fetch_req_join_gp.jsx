"use client"

export default async function FetchJoinToGroup(groupId) {
    try {
        const rep = await fetch("http://localhost:8080/api/RequestToJoinGroups", {
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

        return "pending"
       
    } catch (error) {
        console.error("Error sending the join request:", error.message || error);
    }
}
