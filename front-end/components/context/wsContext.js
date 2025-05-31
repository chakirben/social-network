'use client'
import { usePathname } from 'next/navigation'
import { createContext, useEffect, useRef, useState } from 'react'

export const WebSocketContext = createContext(null)

export const WebSocketProvider = ({ children }) => {
  const [statuses, setStatuses] = useState({})
  const [discussionMap, setDiscussionMap] = useState({})
  // Group5: [
  //   { content: "hello", sentAt: "01-01-2002" },
  //   { content: "hi again", sentAt: "01-01-2002" }
  // ],
  // user1: [
  //   { content: "hey", sentAt: "01-01-2002" }
  // ]
  const [socket, setSocket] = useState(null)
  const connectedRef = useRef(false)
  const pathname = usePathname()

  const Connect = () => {
    if (connectedRef.current) return
    const ws = new WebSocket('ws://localhost:8080/api/ws')
    ws.addEventListener("open", () => {
      console.log("WebSocket connected")
      ws.send(JSON.stringify({ type: "hello" }))
      setSocket(ws)
      connectedRef.current = true
    })

    ws.addEventListener("close", () => {
      console.log("WebSocket disconnected")
      setSocket(null)
      connectedRef.current = false
    })

    ws.addEventListener("error", (e) => {
      console.error("WebSocket error", e)
    })

    ws.addEventListener("message", (event) => {
      try {
        const data = JSON.parse(event.data);
        console.log("Message from server:", data);

        switch (data.type) {
          case "status":
            // handle status update
            console.log("Status update:", data);
            if (data.userId && data.statusType) {
              setStatuses(prev => ({
                ...prev,
                [data.userId]: data.statusType
              }))
            }
            break;
          case "message":
            // handle chat/message
            console.log("New message:", data);
          case "message":
            console.log("New message:", data)

            if (data.scope === "group" && data.groupId) {
              const discussionKey = "group" + data.groupId

              setDiscussionMap(prev => ({
                ...prev,
                [discussionKey]: [...(prev[discussionKey] || []), {
                  content: data.content,
                  sentAt: data.sentAt,
                  senderId: data.senderId
                }]
              }))
            } else if (data.scope === "private" && data.userId) {
              const discussionKey = "user"+data.userId 

              setDiscussionMap(prev => ({
                ...prev,
                [discussionKey]: [...(prev[discussionKey] || []), {
                  content: data.content,
                  sentAt: data.sentAt,
                  senderId: data.senderId
                }]
              }))
            }
            break

          case "notification":
            // handle notification
            console.log("Notification:", data);
            break;

          default:
            console.warn("Unknown message type:", data.type);
        }
      } catch (err) {
        console.error("Invalid JSON from WebSocket:", event.data);
      }
    });

    return ws
  }

  useEffect(() => {
    if (connectedRef.current) return
    if (pathname === "/login" || pathname === "/register") return
    Connect()
  }, [pathname])

  return (
    <WebSocketContext.Provider value={{ socket, Connect, discussionMap, setDiscussionMap, statuses, setStatuses }}>
      {children}
    </WebSocketContext.Provider>
  )
}