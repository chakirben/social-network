'use client'

import { createContext, useContext, useEffect, useRef, useState } from 'react'

const WebSocketContext = createContext(null)

export const WebSocketProvider = ({ children }) => {
  const [socket, setSocket] = useState(null)

  useEffect(() => {
    const ws = new WebSocket('ws://localhost:8080/api/ws')
    setSocket(ws)

    const handleBeforeUnload = () => {
      ws.close()
    }

    window.addEventListener('beforeunload', handleBeforeUnload)

    ws.onopen = () => console.log('Connected')
    ws.onclose = () => console.log('Disconnected')

    return () => {
      window.removeEventListener('beforeunload', handleBeforeUnload)
      ws.close()
    }
  }, [])


  return (
    <WebSocketContext.Provider value={socket}>
      {children}
    </WebSocketContext.Provider>
  )
}

export const useWebSocket = () => useContext(WebSocketContext)
