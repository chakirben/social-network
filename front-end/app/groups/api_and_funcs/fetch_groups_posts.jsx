"use client"
import { useRouter } from 'next/navigation';
import { useEffect, useState } from "react";
import Post from "@/components/post";
import CreatPostInGroup from "@/components/groups/creat_postgroup";
import "../../../styles/global.css";
import GroupEventsPage from '@/components/events/groupEventsPage';
import Divider from '@/components/divider';
import CreateEvent from '@/components/createEventForm';


export default function GroupDetails({ groupId, title }) {
    const [PostsGroup, setPostsGroup] = useState([]);
    const [isLoading, setIsLoading] = useState(true);
    const [events, setEvents] = useState([])
    const [activeTab, setActiveTab] = useState("posts");
    useEffect(() => {
        console.log("Updated events:", events);
    }, [events]);
    const router = useRouter();

    const fetchPosts = async () => {
        setIsLoading(true);
        try {
            const rep = await fetch(`http://localhost:8080/api/PostsGroups?id=${groupId}`, {
                credentials: "include",
                headers: { "Content-Type": "application/json" },
            });
            const PostsGroupData = await rep.json();
            setPostsGroup(PostsGroupData || []);
        } catch (error) {
            console.error("Error fetching posts:", error);
        } finally {
            setIsLoading(false);
        }
    };
    const fetchEvents = async () => {
        try {
            const res = await fetch(`http://localhost:8080/api/GetGroupEvents?id=${groupId}`, {
                credentials: "include",
                headers: { "Content-Type": "application/json" },
            });
            const data = await res.json();
            setEvents(data || []);
        } catch (error) {
            console.error("Error fetching events:", error);
        }
    };

    useEffect(() => {
        if (activeTab === "events") {
            fetchEvents();
        }
    }, [groupId, activeTab]);


    useEffect(() => {
        if (activeTab === "posts") {
            fetchPosts();
        }
    }, [groupId, activeTab]);

    return (
        <div className='postEventsInGroup'>
            <div className="GPTitle">
                <div className="btnback">
                    <img src="./../images/arrow-left.svg" />
                    <button className="backbtn" onClick={() => router.push(`/groups`)}>back</button>
                </div>
                <div className="titleandimg">
                    <img src="./../images/group.svg" />
                    <p>{title}</p>
                </div>
            </div>

            <div className='filterPostsAndEvents'>
                <span
                    className={`postsSpan ${activeTab === "posts" ? "active" : ""}`}
                    onClick={() => setActiveTab("posts")}
                >
                    Posts
                </span>
                <span
                    className={`eventsSpan ${activeTab === "events" ? "active" : ""}`}
                    onClick={() => setActiveTab("events")}
                >
                    Events
                </span>
            </div>
            <Divider />
            {isLoading && <div>Loading...</div>}

            {!isLoading && activeTab === "posts" && (
                <>
                    <CreatPostInGroup gpid={groupId} />
                    {PostsGroup.map((pst) => (
                        <Post key={pst.id} pst={pst} />
                    ))}
                </>
            )}

            {!isLoading && activeTab === "events" && (
                <>
                    <CreateEvent setEvents={setEvents} evnts={events} />
                    <GroupEventsPage id={groupId} events={events} />
                </>
            )}
        </div>
    );
}
