'use client'
import { useEffect, useState } from "react";
import SideBar from "@/components/sidebar";
import "./home.css";
import "../../styles/global.css";
import "./../groups/css/groups1.css"


import Post from "@/components/post";
import SearchBar from "@/components/searchBar";
import CreatePost from "@/components/creatPostForm";
import SearchTerm from "@/components/search_term";
import HomeEvents from "@/components/events/homeEvents";

export default function Home() {
    const [posts, setPosts] = useState([]);
    const [searchTerm, setSearchTerm] = useState("");
    const [searchResults, setSearchResults] = useState(null); // optional if needed

    // fetch all posts initially
    useEffect(() => {
        const fetchPosts = async () => {
            try {
                const response = await fetch("http://localhost:8080/api/GetPosts", {
                    credentials: "include"
                });
                const data = await response.json();
                setPosts(data);
            } catch (error) {
                console.error("Error fetching posts:", error);
            }
        };
        fetchPosts();
    }, []);

    useEffect(() => {
        const timeout = setTimeout(() => {
            if (searchTerm.trim() !== "") {
                FetchSearch(searchTerm);
            }
        }, 300);

        return () => clearTimeout(timeout); // clean up the previous timer
    }, [searchTerm]);


    const handlernewpost = (newpost) => {
        setPosts((prev) => [newpost , ...prev])
    }

    return (
        <div className="home">
            <SideBar />
            <div className="homeP">
                <SearchBar searchTerm={searchTerm} setSearchTerm={setSearchTerm} />

                {searchTerm.trim() === "" ? (
                    <div className="sc">
                        <CreatePost newpost={handlernewpost}/>
                        <HomeEvents />
                        <div className="posts">
                            {posts.length === 0 ? (
                                <div className="loading">Loading posts...</div>
                            ) : (
                                posts.map((post) => (
                                    <Post key={post.id || post._id} pst={post} />
                                ))
                            )}
                        </div>
                    </div>
                ) : (

                    <SearchTerm search={searchTerm} />

                )}
            </div>
        </div>
    );
}
