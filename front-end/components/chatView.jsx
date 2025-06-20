"use client"

import { useContext, useEffect } from "react"
import { usePathname } from "next/navigation"
import { WebSocketContext } from "./context/wsContext"
import MessageInput from "./MessageInput"

export default function ChatView() {
  const { discussionMap, setDiscussionMap } = useContext(WebSocketContext)
  const pathname = usePathname()

  const match = pathname.match(/^\/chat\/(user|group)(\d+)$/)
  const type = match?.[1]
  const id = match?.[2]

  const discussionKey = type + id
  const messages = discussionMap[discussionKey] || []

  useEffect(() => {
    if (!type || !id || discussionMap[discussionKey]) return

    const fetchMessages = async () => {
      try {
        const baseUrl = `${process.env.NEXT_PUBLIC_API_URL}/api/ws/messages`
        const queryParams =
          type === "user"
            ? `?type=private&other_id=${id}`
            : `?type=group&group_id=${id}`

        const res = await fetch(baseUrl + queryParams, {
          credentials: "include",
        })
        if (!res.ok) throw new Error("Failed to fetch messages")
        const data = await res.json()
        setDiscussionMap((prev) => ({
          ...prev,
          [discussionKey]: data,
        }))
      } catch (err) {
        console.error("Error fetching messages:", err)
      }
    }

    fetchMessages()
  }, [type, id, discussionMap, discussionKey, setDiscussionMap])

  return (
    <div className="chatView df cl">
      <div className="MessagesContainer">
        {messages.map((msg, idx) => (
          <div key={idx} className="message">
            <strong>{msg.senderId}</strong>: {msg.content}
          </div>
        ))}
      </div>
      <MessageInput />
    </div>
  )
}
