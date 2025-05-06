'use client'
import { useEffect, useState } from "react";
import SideBar from "@/components/sidebar";
import "./home.css"
import "../../styles/global.css"
import Post from "@/components/post";
import SearchBar from "@/components/searchBar";
import CreatePost from "@/components/creatPostForm";
import Notification from "@/components/notification";
export default function Home() {
    const [posts, setPosts] = useState([]);

    useEffect(() => {
        const fetchPosts = async () => {
            try {
                const response = await fetch("http://localhost:8080/api/GetPosts", { credentials: "include" });
                const data = await response.json();
                setPosts(data);
            } catch (error) {
                console.error("Error fetching posts:", error);
            }
        };

        fetchPosts();
    }, []);

    let notif = {
        id: 1,
        firstName: 'bilalalal',
        notificationType: 'eventRequest',
        image: '/uploads/1745850185765870914_avatar.jpg',
        prjOrEvent: 'Foot ball match',
        time: '2025-05-05 09:56:23',

    }

    return (
        <div className="home">
            <SideBar />
            <div>
                <SearchBar />
                <CreatePost></CreatePost>
                <Notification notification={notif} />
                <div className="posts">
                    {posts.length === 0 ? (
                        <div>Loading posts...</div>
                    ) : (
                        posts.map((post) => (
                            <Post key={post.id || post._id} pst={post} />
                        ))
                    )}
                </div>
            </div>
        </div>
    );
}
