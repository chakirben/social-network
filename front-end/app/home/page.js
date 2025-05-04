'use client'
import { useEffect, useState } from "react";
import SideBar from "@/components/sidebar";
import "./home.css"
import Post from "@/components/post";

export default function Home() {
    // State to hold the fetched posts
    const [posts, setPosts] = useState([]);

    // Fetch posts from the backend
    useEffect(() => {
        const fetchPosts = async () => {
            try {
                const response = await fetch("http://localhost:8080/api/GetPosts", { credentials: "include" }); // Replace with your backend API URL
                const data = await response.json();
                setPosts(data);
            } catch (error) {
                console.error("Error fetching posts:", error);
            }
        };

        fetchPosts();
    }, []);

    return (
        <div className="home">
            <SideBar />
            <div>
                <p>Welcome Home Honey</p>
                <div className="posts">
                    {posts.length === 0 ? (
                        <div>Loading posts...</div>
                    ) : (
                        posts.map((post) => (
                            <Post pst={post} />
                        ))
                    )}
                </div>
            </div>
        </div>
    );
}
