"use client"

import SideBar from "@/components/sidebar"
import { useContext, useEffect, useCallback } from "react"
import { WebSocketContext } from "@/components/context/wsContext.js"
import DiscussionCard from "@/components/context/discussionCard.jsx"
import "./chat.css"

export default function ChatLayout({ children }) {
    const { discussionMap, setDiscussionMap } = useContext(WebSocketContext)
    const { statuses, setStatuses } = useContext(WebSocketContext)

    // Fetch discussions
    useEffect(() => {
        const fetchDiscussions = async () => {
            try {
                const response = await fetch("http://localhost:8080/api/GetDiscussionList", {
                    credentials: "include",
                })
                const data = await response.json()
                const map = {}
                data.forEach((discussion) => {
                    const key = discussion.isGroup ? `group${discussion.id}` : `user${discussion.id}`
                    map[key] = [discussion]
                })
                setDiscussionMap(map)
            } catch (error) {
                console.error("Error fetching discussions:", error)
            }
        }
        fetchDiscussions()
    }, [])

    // Fetch online users
    const fetchOnlineUsers = useCallback(async () => {
        try {
            const response = await fetch("http://localhost:8080/api/online", {
                credentials: "include"
            })
            const users = await response.json()
            const newStatuses = {}
            users.forEach(user => {
                const id = user.id || user._id
                newStatuses[id] = {
                    id,
                    firstName: user.firstName,
                    lastName: user.lastName,
                    avatar: user.avatar,
                    isOnline: true,
                    lastSeen: new Date().toISOString()
                }
            })
            setStatuses(newStatuses)
        } catch (error) {
            console.error("Error fetching online users:", error)
        }
    }, [setStatuses])

    useEffect(() => {
        fetchOnlineUsers()
        const interval = setInterval(fetchOnlineUsers, 10000)
        return () => clearInterval(interval)
    }, [fetchOnlineUsers])

    return (
        <div className="df">
            <SideBar />
            <div className="leftSection df cl">
                <h3 className="Msgs">Online ({statuses ? Object.keys(statuses).length : 0})</h3>
                <div className="online-users-container">
                    {statuses && Object.entries(statuses).map(([userId, user]) => (
                        <div key={userId} className="online-user-avatar" title={`${user.firstName} ${user.lastName}`}>
                            <div className="avatar-wrapper">
                                <img 
                                    src={user.avatar || 'images/Avatars.png'} 
                                    alt={`${user.firstName} ${user.lastName}`}
                                    className="online-avatar"
                                    onError={(e) => { e.target.src = 'images/Avatars.png' }}
                                />
                                <div className="online-indicator"></div>
                            </div>
                        </div>
                    ))}
                    {(!statuses || Object.keys(statuses).length === 0) && (
                        <div className="no-online-users">No users online</div>
                    )}
                </div>

                <h3 className="Msgs">Messages ({discussionMap ? Object.keys(discussionMap).length : 0})</h3>
                <div className="discussionList df cl">
                    {discussionMap && Object.entries(discussionMap).map(([key, messages]) => {
                          if (!Array.isArray(messages) || messages.length === 0) return null
                        const discussion = messages[0]
                        return <DiscussionCard key={key} discussion={discussion} />
                    })}
                    {(!discussionMap || Object.keys(discussionMap).length === 0) && (
                        <div className="no-discussions">No discussions available</div>
                    )}
                </div>
            </div>

            <div style={{ flex: 1 }}>
                {children}
            </div>
        </div>
    )
}
