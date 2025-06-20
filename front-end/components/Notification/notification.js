"use client";

import styles from "./notification.module.css";
import Avatar from "../avatar/avatar";
import { useState } from "react";

export default function Notification({ notification }) {
    const [response, setResponse] = useState(null);

    const {
        id,
        senderFirstName,
        senderLastName,
        senderAvatar,
        type,
        status,
        notificationDate,
    } = notification;

    const name = `${senderFirstName} ${senderLastName}`;

    const getMessage = () => {
        switch (type) {
            case "follow_request":
                return `${name} requested to follow you`;
            case "group_invite":
                return `${name} invited you to join a group`;
            case "group_join_request":
                return `${name} requested to join your group`;
            case "new_event":
                return `${name} created a new event`;
            default:
                return "You have a new notification";
        }
    };

    const handleResponse = async (action) => {
        try {
            const res = await fetch(
                `${process.env.NEXT_PUBLIC_API_URL}/api/respondToNotification?notificationId=${id}&actionType=${action}`,
                {
                    method: "POST",
                    credentials: "include",
                }
            );
            if (res.ok) {
                setResponse(action);
            } else {
                console.error("Failed to respond to notification");
            }
        } catch (err) {
            console.error("Error:", err);
        }
    };

    return (
        <div className={styles.notification}>
            <div className={styles.info}>
                <Avatar url={senderAvatar} name={senderFirstName} size="small" />
                <div className={styles.text}>
                    <p className={styles.message}>{getMessage()}</p>
                    <p className={styles.date}>{new Date(notificationDate).toLocaleString()}</p>
                </div>
            </div>
            <div className={styles.buttons}>
                {response === "accept" && <span className={styles.accepted}>Accepted ! </span>}
                {response === "decline" && <span className={styles.declined}>Declined !</span>}
                {!response && status === "pending" && (type === "follow_request" || type === "group_invite" || type === "group_join_request"  || type === "new_event") && (
                    <>
                        <span className={styles.accept} onClick={() => handleResponse("accept")}>Accept</span>
                        <span className={styles.decline} onClick={() => handleResponse("decline")}>Decline</span>
                    </>
                )}
            </div>
        </div>
    );
}
