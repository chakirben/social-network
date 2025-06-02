'use client';
import { useState } from "react";
import { timePassed } from "@/public/utils/timePassed";
import { useUser } from "../context/userContext";
import styles from './message.module.css';

export default function Message({ msg }) {
    const { user } = useUser();
    const isSender = msg.sender_id === user.id;

    const messageClass = isSender ?  `${styles.message} ${styles.me}` : styles.message ;

    const [showDetails, setShowDetails] = useState(false);

    const handleClick = () => {
        setShowDetails(prev => !prev);
    };

    return (
        <div className={messageClass} onClick={handleClick}>
            <div className={styles.MsgContent}>{msg.content}</div>
            <div className={showDetails ? styles.date : `${styles.date} ${styles.hidden}`}>
                {timePassed(msg.sent_at)}
            </div>
        </div>
    );
}
