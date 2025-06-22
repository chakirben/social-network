'use client';

import { use, useEffect, useState } from "react";
import SideBar from "@/components/sidebar";
import "../../home/home.css";
import "../../../styles/global.css";
import Post from "@/components/post";
import "../../home/comments.css";
import Comment from "@/components/comments";
import Header from "@/components/Header/header";
import CommentInput from "@/components/commentInput";

export default function PostPage({ params }) {
    const [post, setPost] = useState(null);
    const [comments, setComments] = useState([]);
    const [unauthorized, setUnauthorized] = useState(false);
    const { id } = use(params);
    useEffect(() => {
        const fetchPost = async () => {
            try {
                const response = await fetch(`/api/GetOnePost?id=${id}`, {
                    credentials: "include",
                });

                if (response.status === 401) {
                    setUnauthorized(true);
                    return;
                }

                const data = await response.json();
                setPost(data);
            } catch (error) {
                console.error("Error fetching post:", error);
            }
        };

        fetchPost();
    }, [id]);

    useEffect(() => {
        const fetchComments = async () => {
            try {
                const response = await fetch(`/api/GetComments?id=${id}`, {
                    credentials: "include",
                });

                if (response.status === 401) {
                    return;
                }

                const data = await response.json();
                setComments(data);
            } catch (error) {
                console.error("Error fetching comments:", error);
            }
        };

        fetchComments();
    }, [id]);

    return (
        <div className="home">
            <SideBar />
            <div className="homeP df cl spB">
                <Header pageName={"post"} />
                <div className="postContent df cl">
                    {unauthorized ? (
                        <div className="loading">Unauthorized</div>
                    ) : post === null ? (
                        <div className="loading">Loading post...</div>
                    ) : (
                        <>
                            <Post key={post.id || post._id} pst={post} />

                            {comments && comments.length === 0 ? (
                                <div className="loading">Loading comments...</div>
                            ) : (
                                <div className="commentsSection">
                                    <h4 className="cmsss">Comments</h4>
                                    {comments?.map((cmnt) => (
                                        <Comment key={cmnt.id || cmnt._id} comment={cmnt} />
                                    ))}
                                </div>
                            )}

                            <CommentInput id={id} setComments={setComments} />
                        </>
                    )}
                </div>
            </div>
        </div>
    );

}