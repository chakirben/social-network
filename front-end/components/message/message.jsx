'use client';
import { useState } from "react";
import { timePassed } from "@/public/utils/timePassed";
import { useUser } from "../context/userContext";
import styles from './message.module.css';
import Avatar from "../avatar/avatar";

export default function Message({ msg }) {
    const { user } = useUser();
    const isSender = msg.sender_id === user?.id;

    const messageContainerClass = isSender
        ? `${styles.messageContainer} ${styles.me}`
        : styles.messageContainer;
    const messageClass = isSender
        ? `${styles.message} ${styles.me}`
        : styles.message;
    const [showDetails, setShowDetails] = useState(false);

    const handleClick = () => {
        setShowDetails(prev => !prev);
    };

    return (
        <div className={messageContainerClass}>
            {!isSender && (
               <Avatar name={msg.first_name} url={msg.avatar}/>
            )}

            <div onClick={handleClick} className={messageClass}>
                <div className={styles.chatName}>{isSender ? "" : `${msg.first_name}`}</div>
                <div className={styles.MsgContent}>{msg.content}</div>
                <div className={showDetails ? styles.date : `${styles.date} ${styles.hidden}`}>
                    {msg.sent_at ? timePassed(msg.sent_at) : "now"}
                </div>
            </div>
        </div>
    );
}
