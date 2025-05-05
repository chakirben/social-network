"use client";
import Groups from "@/components/groups";
import SideBar from "@/components/sidebar";
import { useEffect, useState } from "react";
export default function GroupsPage() {
    const [GroupsData, setGroups] = useState([])

    useEffect(() => {
        const fetchGroups = async () => {
            try {
                const response = await fetch("http://localhost:8080/api/NotMyGroups", {
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
                   <SideBar />
                <div>   
                    <div className="groups">
                        <div>
                            {Groups.length === 0 ? (
                                <div>Loading groups...</div>
                            ) : (
                                GroupsData.map((Group) => (
                                    <Groups grp={Group}/>
                                ))
                            )}
                        </div>
                    </div>
                </div>
            </div>
    );
}