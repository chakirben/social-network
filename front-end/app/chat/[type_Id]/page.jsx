'use client'

import { useContext, useEffect, useState } from "react"
import { usePathname } from "next/navigation"
import MessageInput from "@/components/chatInput/chatInput"
import Message from "@/components/message/message"
import Header from "@/components/Header/header"
import { WebSocketContext } from "@/components/context/wsContext"
import { useUser } from "@/components/context/userContext"

export default function ChatView() {
  const pathname = usePathname()
  const [messages, setMessages] = useState([])
  const { wsMessages } = useContext(WebSocketContext)
  const { myId } = useUser()

  const match = pathname.match(/^\/chat\/(user|group)(\d+)_([\w_]+)$/)
  const type = match?.[1]
  const id = match?.[2]
  const rawName = match?.[3] || ''
  const name = rawName.replace(/_/g, ' ')

  useEffect(() => {
    if (!type || !id) return
    const fetchMessages = async () => {
      try {
        const baseUrl = `${process.env.NEXT_PUBLIC_API_URL}/api/fetchMessages`

        const queryParams =
          type === "user"
            ? `?type=private&other_id=${id}`
            : `?type=group&group_id=${id}`

        const res = await fetch(baseUrl + queryParams, {
          credentials: "include",
        })

        if (!res.ok) throw new Error("Failed to fetch messages")

        const data = await res.json()
        setMessages(data || [])
      } catch (err) {
        console.error("Error fetching messages:", err)
      }
    }
    fetchMessages()
  }, [type, id, pathname])

  useEffect(() => {
    if (wsMessages.length === 0) return;

    const lastMsg = wsMessages[wsMessages.length - 1];
    console.log("last msg is " ,lastMsg);
    
    const isRelevant =
      (type === "user" && (lastMsg.receiver_id == id || lastMsg.sender_id == id)) ||
      (type === "group" && lastMsg.groupId == id);

    if (isRelevant) {
      setMessages((prev) => [...prev, lastMsg]);
    }
  }, [wsMessages, id, type]);

  return (
    <div className="chatView df cl">
      <Header pageName={name || 'chat'} />
      <div className="MessagesContainer">
        {messages.map((msg, idx) => (
          <Message msg={msg} key={idx} />
        ))}
      </div>
      <MessageInput id={id} type={type} />
    </div>
  )
}