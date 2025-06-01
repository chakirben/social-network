'use client'

import { usePathname } from 'next/navigation'
import { createContext, useEffect, useRef, useState } from 'react'

export const WebSocketContext = createContext(null)

export const WebSocketProvider = ({ children }) => {
  const [statuses, setStatuses] = useState({})
  const [discussionMap, setDiscussionMap] = useState({})
  const [socket, setSocket] = useState(null)

  const connectedRef = useRef(false)
  const pathname = usePathname()

  const Connect = () => {
    if (connectedRef.current) return

    const ws = new WebSocket('ws://localhost:8080/api/ws')

    ws.addEventListener('open', () => {
      console.log('WebSocket connected')
      ws.send(JSON.stringify({ type: 'hello' }))
      setSocket(ws)
      connectedRef.current = true
    })

    ws.addEventListener('close', () => {
      console.log('WebSocket disconnected')
      setSocket(null)
      connectedRef.current = false
    })

    ws.addEventListener('error', (e) => {
      console.error('WebSocket error', e)
    })

    ws.addEventListener('message', (event) => {
      try {
        const data = JSON.parse(event.data)
        console.log('Message from server:', data)

        switch (data.type) {
          case 'Status':
            // live status update with full user data
            const { userId, statusType, user } = data

            if (userId && statusType) {
              setStatuses((prev) => {
                const updated = { ...prev }

                if (statusType === 'online') {
                  if (user) {
                    updated[userId] = {
                      firstName: user.firstName,
                      lastName: user.lastName,
                      avatar: user.avatar
                    }
                  }
                } else if (statusType === 'offline') {
                  delete updated[userId]
                }

                return updated
              })
            }

            break

          case 'message':
          case 'grpMessage':
            const discussionKey =
              data.scope === 'group'
                ? 'group' + data.groupId
                : 'user' + data.userId

            setDiscussionMap((prev) => ({
              ...prev,
              [discussionKey]: [
                ...(prev[discussionKey] || []),
                {
                  content: data.content,
                  sentAt: data.sentAt,
                  senderId: data.senderId
                }
              ]
            }))
            break

          case 'notification':
            console.log('Notification:', data)
            break

          default:
            console.warn('Unknown message type:', data.type)
        }
      } catch (err) {
        console.error('Invalid JSON from WebSocket:', event.data)
      }
    })

    return ws
  }

  useEffect(() => {
    if (connectedRef.current) return
    if (pathname === '/login' || pathname === '/register') return
    Connect()
  }, [pathname])

  return (
    <WebSocketContext.Provider
      value={{
        socket,
        Connect,
        discussionMap,
        setDiscussionMap,
        statuses,
        setStatuses
      }}
    >
      {children}
    </WebSocketContext.Provider>
  )
}