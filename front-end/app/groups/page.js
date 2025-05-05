"use client";
import Goups from "@/components/groups";
import SideBar from "@/components/sidebar";
import { useEffect, useState } from "react";
export default function GroupsPage() {
    const [Groups, setGroups] = useState([])

    useEffect(() => {
        const fetchGroups = async () => {
            try {
                const response = await fetch("http://localhost:8080/api/MyGroups", {
                    credentials:"include",
                } ); 
                const data = await response.json();
                console.log(data)
                setGroups(data);
            } catch (error) {
                console.error("Error fetching groups:", error);
            }
        }
        fetchGroups()
    }, [])
     return (
            <div className="home">
                <div>
                    <div>
                        {Groups.length === 0 ? (
                            <div>Loading groups...</div>
                        ) : (
                            Groups.map((Group) => (
                                <Goups grp={Group}/>
                            ))
                        )}
                    </div>
                </div>
            </div>
    );
}