'use client'
import { useEffect, useState } from "react";
import SideBar from "@/components/sidebar";
import "./home.css"
import "../../styles/global.css"
import Post from "@/components/post";
import SearchBar from "@/components/searchBar";
import CreatePost from "@/components/creatPostForm";
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

    return (
        <div className="home">
            <SideBar />
            <div>
                <SearchBar />
                <CreatePost></CreatePost>
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
