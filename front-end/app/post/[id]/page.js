'use client'

import { use, useEffect, useState } from "react";
import SideBar from "@/components/sidebar";
import "../../home/home.css"
import "../../../styles/global.css"
import Post from "@/components/post";
import "../../home/comments.css"
import Comment from "@/components/comments";
import Header from "@/components/Header/header";
import CommentInput from "@/components/commentInput";

export default function PostPage({ params }) {
    const [post, setPost] = useState(null);  // Use null initially for the post
    const [comments, setComments] = useState([]);

    const { id } = use(params);  // No need for 'use', just use params directly

    // Fetch post data
    useEffect(() => {
        const fetchPost = async () => {
            try {
                const response = await fetch(`http://localhost:8080/api/GetOnePost?id=${id}`, { credentials: "include" });
                const data = await response.json();
                setPost(data);
            } catch (error) {
                console.error("Error fetching post:", error);
            }
        };
        fetchPost();
    }, [id]); // Make sure to use id as a dependency

    // Fetch comments data
    useEffect(() => {
        const fetchComments = async () => {
            try {
                const response = await fetch(`http://localhost:8080/api/GetComments?id=${id}`, { credentials: "include" });
                const data = await response.json();
                setComments(data);
            } catch (error) {
                console.error("Error fetching comments:", error);
            }
        };
        fetchComments();
    }, [id]); // Use id as a dependency to fetch comments when id changes

    return (
        <div className="home">
            <SideBar />
            <div className="singlePost df cl spB">
                <Header pageName={"post"}></Header>
                <div className="postContent df cl">
                    {post === null ? (
                        <div className="loading">Loading post...</div>
                    ) : (
                        <Post key={post.id || post._id} pst={post} />
                    )}
                    {comments && comments.length === 0 ? (
                        <div className="loading">Loading comments...</div>
                    ) : (
                        comments && comments.map((cmnt) => (
                            <Comment key={cmnt.id || cmnt._id} comment={cmnt} />
                        ))
                    )}
                </div>
                <CommentInput id={id} setComments={setComments}></CommentInput>
            </div>
        </div>
    );
}