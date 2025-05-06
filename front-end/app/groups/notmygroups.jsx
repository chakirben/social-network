"use client";
import NoMyGroup from "@/components/notmygroups";
import { useEffect, useState } from "react";
import "./groups.css"

export default function NotMyGroupsPage() {
    const [NotMygroupsData, setGroups] = useState([]);

    useEffect(() => {
        const fetchGroups = async () => {
            try {
                const response = await fetch("http://localhost:8080/api/NotMyGroups", { credentials: "include" });
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
            {NotMygroupsData.length === 0 ? (
                <div>Loading groups...</div>
            ) : (
                NotMygroupsData.map((group) => (
                    <NoMyGroup key={group.Id || ''} group={group} />
                ))
            )}
        </div>

    )
}

