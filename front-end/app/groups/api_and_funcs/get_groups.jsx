"use client";
import { useEffect, useState } from "react";
import MyGroup from "@/components/groups/mygroups";
import NoMyGroup from "@/components/groups/notmygroups";

export default function MyGroupsPage({ onJoin, onView }) {
    const [myGroups, setMyGroups] = useState([]);
    const [notMyGroups, setNotMyGroups] = useState([]);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        const fetchGroups = async () => {
            try {
                const [myGroupsRes, notMyGroupsRes] = await Promise.all([
                    fetch("http://localhost:8080/api/MyGroups", { credentials: "include"  }),
                    fetch("http://localhost:8080/api/NotMyGroups", { credentials: "include" }),
                ]);
                const myGroupsData = await myGroupsRes.json();
                const notMyGroupsData = await notMyGroupsRes.json();
                setMyGroups(myGroupsData || []);
                setNotMyGroups(notMyGroupsData || []);
            } catch (error) {
                console.error("Error fetching groups:", error);
            } finally {
                setIsLoading(false);
            }
        };

        fetchGroups();
    }, []);

    if (isLoading) {
        return <div>Loading groups...</div>;
    }

    if (myGroups.length === 0 && notMyGroups.length === 0) {
        return <div>No groups found...</div>;
    }

    return (
        <div className="groupsmn">
            {myGroups.map((group) => (
                <MyGroup key={group.Id} group={group} onView={onView} />
            ))}
            {notMyGroups.map((group) => (
                <NoMyGroup key={group.Id} group={group} onJoin={onJoin} />
            ))}
        </div>
    );
}
