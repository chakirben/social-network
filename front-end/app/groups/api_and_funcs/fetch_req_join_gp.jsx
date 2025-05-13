"use client"

export default async function FetchJoinToGroup(groupId) {
    try {
        const rep = await fetch("http://localhost:8080/api/RequestToJoinGroups", {
            method: "POST",
            credentials: "include",
            headers: {
                "Content-Type": "application/json"
            },
            body : JSON.stringify({
                groupId : groupId
            })
        
        })
        if (!rep.ok) {
            throw new Error("Failed to send Request to join the group");
        }
        const repData = await rep.json();
        console.log("Request send successfully:", repData);        
    } catch {
        console.error("Error creating the group:", error);
    }
}