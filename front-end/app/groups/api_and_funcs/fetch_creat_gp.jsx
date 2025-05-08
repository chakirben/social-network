"use client"
export default async function FetchCreatGroup(title, description) {
    try {
        const rep = await fetch("http://localhost:8080/api/CreatGroup", {
            method: "POST",
            credentials: "include",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                title: title,
                description: description,
            }),
        });

        if (!rep.ok) {
            throw new Error("Failed to create the group");
        }

        const repData = await rep.json();
        console.log("Group created successfully:", repData);
    } catch (error) {
        console.error("Error creating the group:", error);
        return null;
    }
}
