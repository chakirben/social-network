'use client';

import { usePathname } from 'next/navigation';
import { createContext, useEffect, useRef, useState } from 'react';

export const WebSocketContext = createContext(null);

export const WebSocketProvider = ({ children }) => {
  const [statuses, setStatuses] = useState({});
  const [discussionMap, setDiscussionMap] = useState({});
  const [wsMessages, setwsMessages] = useState([]);
  const [socket, setSocket] = useState(null);

  const connectedRef = useRef(false);
  const pathname = usePathname();

  const Connect = () => {
    if (connectedRef.current) return;

    const ws = new WebSocket('ws://localhost:8080/api/ws');

    ws.addEventListener('open', () => {
      console.log('WebSocket connected');
      ws.send(JSON.stringify({ type: 'hello' }));
      setSocket(ws);
      connectedRef.current = true;
    });

    ws.addEventListener('close', () => {
      console.log('WebSocket disconnected');
      setSocket(null);
      connectedRef.current = false;
    });

    ws.addEventListener('error', (e) => {
      console.error('WebSocket error', e);
    });

    ws.addEventListener('message', (event) => {
      try {
        const data = JSON.parse(event.data);
        console.log('Message from server:', data);

        switch (data.type) {
          case 'Status': {
            const { userId, statusType, user } = data;
            if (userId && statusType) {
              setStatuses((prev) => {
                const updated = { ...prev };
                if (statusType === 'online' && user) {
                    updated[userId] = {
                    firstName: user.firstName,
                    lastName: user.lastName,
                    avatar: user.avatar
                  };
                } else if (statusType === 'offline') {
                  delete updated[userId];
                }
                return updated;
              });
            }
            break;
          }
          case 'message': {
            const formattedMsg = {
              content: data.content,
              sender_id: data.sender,
              receiver_id: data.receiverId,
              sent_at: data.sentAt || data.sent_at,
              type: 'private',
            };
            setwsMessages((prev) => [...prev, formattedMsg]);
            break;
          }

          case 'groupmsg': {
            const formattedGroupMsg = {
              content: data.content,
              sender_id: data.sender,
              groupId: data.groupID || data.groupId,
              sent_at: data.sentAt || data.sent_at,
              type: 'group',
            };
            const discussionKey = 'group' + formattedGroupMsg.groupId;
            setDiscussionMap((prev) => ({
              ...prev,
              [discussionKey]: [...(prev[discussionKey] || []), formattedGroupMsg]
            }))
            setwsMessages((prev) => [...prev, formattedGroupMsg]);
            break;
          }

          case 'notification':
            console.log('Notification:', data);
            break;

          default:
            console.warn('Unknown message type:', data.type);
        }
      } catch (err) {
        console.error('Invalid JSON from WebSocket:', event.data);
      }
    });

    return ws;
  };

  useEffect(() => {
    if (connectedRef.current) return;
    if (pathname === '/login' || pathname === '/register') return;
    Connect();
  }, [pathname]);

  return (
    <WebSocketContext.Provider
      value={{
        socket,
        Connect,
        discussionMap,
        setDiscussionMap,
        statuses,
        setStatuses,
        wsMessages,
        setwsMessages,
      }}
    >
      {children}
    </WebSocketContext.Provider>
  );
};