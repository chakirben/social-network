"use client";
import Group from "@/components/groups";
import Groupbar from "@/components/groupbar"
import SideBar from "@/components/sidebar";
import { useEffect, useState } from "react";
import "./groups.css"

export default function GroupsPage() {
    const [groupsData, setGroups] = useState([]);

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
        <div className="home">
            <SideBar />
            <div className="groups">
            <Groupbar />
                {groupsData.length === 0 ? (
                    <div>Loading groups...</div>
                ) : (
                    groupsData.map((group) => (
                        <Group key={group.Id || ''} group={group} />
                    ))
                )}
            </div>
        </div>
    );
}
