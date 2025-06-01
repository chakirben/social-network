'use client'
import { useContext, useEffect } from "react"
import { usePathname } from "next/navigation"
import { WebSocketContext } from "@/components/context/wsContext"
import MessageInput from "@/components/chatInput"

export default function ChatView() {
  const { discussionMap, setDiscussionMap } = useContext(WebSocketContext)
  const pathname = usePathname()

  const match = pathname.match(/^\/chat\/(user|group)(\d+)$/)
  const type = match?.[1]
  const id = match?.[2]

  const discussionKey = type + id
  console.log("Current discussionKey:", discussionKey)

  useEffect(() => {
    if (!type || !id || discussionMap[discussionKey]) return

    const fetchMessages = async () => {
      try {
        const baseUrl = "http://localhost:8080/api/fetchMessages"
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
  }, [type, id, discussionKey, discussionMap, setDiscussionMap])

  // Add this effect to log discussionMap whenever it changes
  useEffect(() => {
    console.log("Updated discussionMap:", discussionMap)
  }, [discussionMap])

  return (
    <div className="chatView df cl">
      <div className="MessagesContainer">
        {(discussionMap[discussionKey]?.messages || []).map((msg, idx) => (
          <div key={idx} className="message">
            <strong>{msg.sender_id}</strong>: {msg.content}
            <div className="timestamp">
              {new Date(msg.sent_at).toLocaleString()}
            </div>
          </div>
        ))}
      </div>
      <MessageInput />
    </div>
  )
}