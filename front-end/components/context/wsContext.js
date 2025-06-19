'use client';

import { usePathname } from 'next/navigation';
import { createContext, useEffect, useRef, useState } from 'react';
import { usePopup } from './popUp';
import { useUser } from './userContext';

export const WebSocketContext = createContext(null);

export const WebSocketProvider = ({ children }) => {
  console.log('WebSocketProvider rendered');

  const [statuses, setStatuses] = useState({});
  const [discussionMap, setDiscussionMap] = useState({});
  const [wsMessages, setwsMessages] = useState([]);
  const [socket, setSocket] = useState(null);
  const { showPopup } = usePopup();
  const connectedRef = useRef(false);
  const pathname = usePathname();
  const [notifCounter, setNotifCounter] = useState(0);
  const [messagesCounter, setMessagesCounter] = useState(0);
  const { user } = useUser();
  const userRef = useRef(user);
  useEffect(() => {
    if (pathname === '/chat') {
      setMessagesCounter(0);
    }
  }, [pathname]);
  useEffect(() => {
    if (pathname === '/notifications') {
      setNotifCounter(0);
    }
  }, [pathname]);
  useEffect(() => {
    userRef.current = user;
  }, [user]);

  const Connect = () => {
    if (connectedRef.current || !userRef.current) return;

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
        const currentUser = userRef.current;
        console.log('Message from server:', data, currentUser);

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
            console.log('Received private message:', data, currentUser);

            if (data.sender !== currentUser?.id) {
              setMessagesCounter((prev) => prev + 1);
            }
            const formattedMsg = {
              content: data.content,
              sender_id: data.sender,
              first_name: data.firstName,
              receiver_id: data.receiverId,
              sent_at: data.sentAt || data.sent_at,
              type: 'private',
            };
            setwsMessages((prev) => [...prev, formattedMsg]);
            break;
          }

          case 'groupmsg': {
            setMessagesCounter((prev) => prev + 1);
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

          case 'Notification':
            console.log(notifCounter);
            
            setNotifCounter((prev) => prev + 1);
            console.log(notifCounter);

            showPopup({ data })
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
    if (!user || connectedRef.current) return;
    if (pathname === '/login' || pathname === '/register') return;
    Connect();
  }, [pathname, user]);

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
        notifCounter,
        setNotifCounter,
        messagesCounter,
        setMessagesCounter
      }}
    >
      {children}
    </WebSocketContext.Provider>
  );
};