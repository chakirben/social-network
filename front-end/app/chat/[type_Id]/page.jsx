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
  const match = pathname.match(/^\/chat\/(user|group)(\d+)$/)
  const type = match?.[1]
  const id = match?.[2]

  // Fetch initial messages
  useEffect(() => {
    if (!type || !id) return

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

    const isRelevant =
      (type === "user" && (lastMsg.sender == id || lastMsg.receiver == id)) ||
      (type === "group" && lastMsg.groupId == id);

    if (isRelevant) {
      // Convert lastMsg to fetched message format
      const formattedMsg = {
        content: lastMsg.content,
        sent_at: lastMsg.sentAt || lastMsg.sent_at, // adjust if key differs
        sender_id: lastMsg.sender,
        receiver_id: lastMsg.receiver,
        groupId: lastMsg.groupId, // if group message, keep groupId
        id: lastMsg.id // if you have an id for deduplication
      };

      setMessages((prev) => {
        return [...prev, formattedMsg];
      });
    }
  }, [wsMessages, id, type]);

  return (
    <div className="chatView df cl">
      <Header pageName={'private chat'} />
      <div className="MessagesContainer">
        {messages.map((msg, idx) => (
          <Message msg={msg} key={idx} />
        ))}
      </div>
      <MessageInput id={id} type={type} />
    </div>
  )
}
