"use client";
import MyGroup from "@/components/mygroups";
import { useEffect, useState } from "react";
import "./groups.css"
export default function MyGroupsPage() {
    const [groupsData, setGroups] = useState([]);

    useEffect(() => {
        const fetchGroups = async () => {
            try {
                const response = await fetch("http://localhost:8080/api/MyGroups", { credentials: "include" });
                const data = await response.json();
                console.log("Fetched groups:", data);
                setGroups(data);
            } catch (error) {
                console.error("Error fetching groups:", error);
            }
        };
        fetchGroups();
    }, []);

    return (
        <div className="groups">
            {groupsData.length === 0 ? (
                <div>Loading groups...</div>
            ) : (
                groupsData.map((group) => (
                    <MyGroup key={group.Id || ''} group={group} />
                ))
            )}
        </div>

    )
}

