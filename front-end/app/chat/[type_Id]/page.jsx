'use client'

import { useContext, useEffect, useState, useRef } from "react"
import { usePathname } from "next/navigation"
import MessageInput from "@/components/chatInput/chatInput"
import Message from "@/components/message/message"
import Header from "@/components/Header/header"
import { WebSocketContext } from "@/components/context/wsContext"
import { useUser } from "@/components/context/userContext"

export default function ChatView() {
  const pathname = usePathname()
  const [messages, setMessages] = useState([])
  const [offset, setOffset] = useState(0)
  const [hasMore, setHasMore] = useState(true)
  const [loading, setLoading] = useState(false)
  const messagesContainerRef = useRef(null)

  const { wsMessages } = useContext(WebSocketContext)
  const { myId } = useUser()

  const match = pathname.match(/^\/chat\/(user|group)(\d+)_([\w_]+)$/)
  const type = match?.[1]
  const id = match?.[2]
  const rawName = match?.[3] || ''
  const name = rawName.replace(/_/g, ' ')

  const fetchMessages = async (newOffset) => {
    if (!type || !id || !hasMore) return;

    setLoading(true);
    try {
      const baseUrl = `/api/fetchMessages`
      const queryParams =
        type === "user"
          ? `?type=private&other_id=${id}&offset=${newOffset}`
          : `?type=group&group_id=${id}&offset=${newOffset}`

      const res = await fetch(baseUrl + queryParams, {
        credentials: "include",
      })

      if (!res.ok) throw new Error("Failed to fetch messages")
      const data = await res.json() || []
      if (data?.length < 10) {
        setHasMore(false)
      }

      if (newOffset === 0) {
        setMessages(data)
        setTimeout(() => {
          messagesContainerRef.current?.scrollTo(0, messagesContainerRef.current.scrollHeight);
        }, 0);
      } else {
        const container = messagesContainerRef.current
        const prevScrollHeight = container.scrollHeight
        setMessages((prev) => [...data, ...prev])
        setTimeout(() => {
          const newScrollHeight = container.scrollHeight
          container.scrollTop = newScrollHeight - prevScrollHeight
        }, 0);
      }
    } catch (err) {
      console.error("Error fetching messages:", err)
    } finally {
      setLoading(false)
    }
  }


  useEffect(() => {
    setOffset(0)
    setHasMore(true)
    fetchMessages(0)
  }, [type, id, pathname])

  useEffect(() => {
    if (!wsMessages.length) return;

    const lastMsg = wsMessages[wsMessages.length - 1];
    const isRelevant =
      (type === "user" && (lastMsg?.receiver_id == id || lastMsg?.sender_id == id)) ||
      (type === "group" && lastMsg?.groupId == id);

    if (isRelevant) {
      setMessages((prev) => Array.isArray(prev) ? [...prev, lastMsg] : [lastMsg]);
      setTimeout(() => {
        messagesContainerRef.current?.scrollTo(0, messagesContainerRef.current.scrollHeight);
      }, 0);
      setOffset((prev) => prev + 1);
    }
  }, [wsMessages, id, type]);

  const handleScroll = (e) => {
    if (!hasMore || loading) return;
    if (e.target.scrollTop === 0) {
      const newOffset = offset + 14
      setOffset(newOffset)
      fetchMessages(newOffset)
    }
  }

  return (
    <div className="chatView df cl">
      <Header pageName={name || 'chat'} />
      <div
        className="MessagesContainer"
        ref={messagesContainerRef}
        onScroll={handleScroll}
        style={{ overflowY: "auto", flex: 1 }}
      >
        {messages?.map((msg, idx) => (
          <Message msg={msg} key={idx} />
        ))}
      </div>
      
      <MessageInput id={id} type={type} />
    </div>
  )
}
