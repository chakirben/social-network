'use client'

import { useEffect, useState, useContext } from "react"
import { useRouter } from "next/navigation"

import SideBar from "@/components/sidebar"
import Post from "@/components/post"
import SearchBar from "@/components/searchBar"
import CreatePost from "@/components/creatPostForm"
import SearchTerm from "@/components/search_term"
import HomeEvents from "@/components/events/homeEvents"
import ProfileButton from "@/components/profileButton"
import UserCard from "@/components/onlineusercard"

import "./home.css"
import "../../styles/global.css"
import "../groups/css/groups1.css"
import { WebSocketContext } from "@/components/context/wsContext"

export default function Home() {
    const router = useRouter()
    const [posts, setPosts] = useState([])
    const [searchTerm, setSearchTerm] = useState("")
    const { statuses, setStatuses } = useContext(WebSocketContext)

    // Fetch posts
    useEffect(() => {
        const fetchPosts = async () => {
            try {
                const response = await fetch("http://localhost:8080/api/GetPosts", {
                    credentials: "include"
                })
                const data = await response.json()
                setPosts(data)
            } catch (error) {
                console.error("Error fetching posts:", error)
            }
        }

        fetchPosts()
    }, [])
    // Fetch online users and set their full info in statuses
    useEffect(() => {
        const fetchOnlineUsers = async () => {
            try {
                const response = await fetch("http://localhost:8080/api/online", {
                    credentials: "include"
                })
                const users = await response.json()

                const newStatuses = {}
                users.forEach(user => {
                    const id = user.id || user._id
                    newStatuses[id] = {
                        firstName: user.firstName,
                        lastName: user.lastName,
                        avatar: user.avatar
                    }
                })

                setStatuses(newStatuses)
            } catch (error) {
                console.error("Error fetching online users:", error)
            }
        }

        fetchOnlineUsers()
        const interval = setInterval(fetchOnlineUsers, 10000)
        return () => clearInterval(interval)
    }, [setStatuses])

    const handleNewPost = (newPost) => {
        setPosts(prev => [newPost, ...prev])
    }

    return (
        <div className="home">
            <SideBar />
            <div className="homeP">
                <SearchBar searchTerm={searchTerm} setSearchTerm={setSearchTerm} />
                {searchTerm.trim() === "" ? (
                    <div className="sc">
                        <CreatePost newpost={handleNewPost} />
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
            <div>
                <ProfileButton />
                <div className="onlineUsers df cl">
                    <h4>Online users</h4>
                    {Object.keys(statuses).length === 0 ? (
                        <div className="loading2">No active users</div>
                    ) : (
                        Object.entries(statuses).map(([uid, userInfo]) => (
                            <UserCard
                                key={uid}
                                user={{
                                    id: uid,
                                    ...userInfo,
                                    status: "online"
                                }}
                            />
                        ))
                    )}
                </div>
            </div>
        </div>
    )
}
