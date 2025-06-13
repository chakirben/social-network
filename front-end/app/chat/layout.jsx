"use client"

import SideBar from "@/components/sidebar"
import { useContext, useEffect, useCallback, useState } from "react"
import { WebSocketContext } from "@/components/context/wsContext"
import DiscussionCard from "@/components/context/discussionCard"
import "./chat.css"
import Avatar from "@/components/avatar/avatar"

export default function ChatLayout({ children }) {
    const { discussionMap, setDiscussionMap } = useContext(WebSocketContext)
    const { statuses, setStatuses } = useContext(WebSocketContext)

    const [friends, setFriends] = useState([])
    const [groups, setGroups] = useState([])

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

    const fetchFriendsAndGroups = useCallback(async () => {
        try {
            const response = await fetch("http://localhost:8080/api/friendsAndGroups", {
                credentials: "include"
            })
            const data = await response.json()
            setFriends(data.friends || [])
            setGroups(data.groups || [])
        } catch (error) {
            console.error("Error fetching friends and groups:", error)
        }
    }, [])

    useEffect(() => {
        fetchFriendsAndGroups()
    }, [fetchFriendsAndGroups])

    return (
        <div className="df">
            <SideBar />
            <div className="leftSection df cl">
                <h3 className="Msgs">People & Groups</h3>
                <div className="online-users-container">
                    {friends.concat(groups).length > 0 ? (
                        friends.concat(groups).map(entity => {
                            const isGroup = !!entity.name
                            const isOnline = !isGroup && statuses[entity.id]?.isOnline
                            const displayName = isGroup
                                ? entity.name
                                : `${entity.firstName} ${entity.lastName}`

                            return (
                                <div
                                    key={entity.id}
                                    className="online-user-avatar"
                                    title={displayName}
                                >
                                   
                                        {entity.avatar ? (
                                            <img
                                                src={entity.avatar}
                                                alt={displayName}
                                                className="online-avatar"
                                            />
                                        ) : (
                                            <Avatar size={48} name={displayName} />
                                        )}
                                        <div
                                            className={`online-indicator ${
                                                isGroup
                                                    ? "group-indicator"
                                                    : isOnline
                                                    ? "online"
                                                    : "offline"
                                            }`}
                                        ></div>
                                    
                                </div>
                            )
                        })
                    ) : (
                        <div className="no-online-users">No users or groups to show</div>
                    )}
                </div>

                <h3 className="Msgs">Messages</h3>
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

            {children}
        </div>
    )
}
