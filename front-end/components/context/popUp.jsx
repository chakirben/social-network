"use client";

import { createContext, useContext, useState, useEffect } from "react";
import styles from "./popup.module.css";
import Image from "next/image";
import { useRouter } from "next/navigation";
import Avatar from "../avatar/avatar";

const PopupContext = createContext();

export function usePopup() {
    return useContext(PopupContext);
}

export function PopupProvider({ children }) {
    console.log("PopupProvider rendered");
    
    const [popup, setPopup] = useState(null);
    const [response, setResponse] = useState(null); 
    const router = useRouter();

    useEffect(() => {
        let timer;
        if (popup) {
            timer = setTimeout(() => {
                setPopup(null);
                setResponse(null); 
            }, 1000);
        }
        return () => clearTimeout(timer);
    }, [popup]);

    const showPopup = (notification) => {
        setPopup(notification.data);
        setResponse(null);
    };

    const handleClick = () => {
        router.push("/notifications");
    };

    const getMessage = (notif) => {
        const name = `${notif.firstName} ${notif.lastName}`;
        switch (notif.notificationType) {
            case "event":
                return (
                    <p>
                        <span className={styles.user}>{name}</span> {"Created Event " + notif.title}
                    </p>
                );
            case "follow":
                return (
                    <p>
                        <span className={styles.user}>{name}</span> Requested to follow you
                    </p>
                );
            case "invite":
                return (
                    <p>
                        <span className={styles.user}>{name}</span> {"Invited you to Join Group " + notif.title}
                    </p>
                );
            case "groupJoinRequest":
                return (
                    <p>
                        <span className={styles.user}>{name}</span> Requested to enter your Group {notif.title}
                    </p>
                );
            default:
                return <p>New notification</p>;
        }
    };

    const handleResponse = async (action) => {
        try {
            
            const res = await fetch(
                `http://localhost:8080/api/respondToNotification?notificationId=${popup.id}&actionType=${action}`,
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
        <PopupContext.Provider value={{ showPopup }}>
            {children}
            {popup && (
                <div className={styles.popupContainer}>
                    <div className={styles.popup}>
                        <div className={styles.info} onClick={handleClick}>
                            <Avatar url={popup.avatar} name={popup.firstName} size="small" />
                            <div className={styles.text}>{getMessage(popup)}</div>
                        </div>

                        <div className={styles.buttons}>
                            {response === "accepted" && <span className={styles.accept}>Accepted</span>}
                            {response === "declined" && <span className={styles.decline}>Declined</span>}

                            {!response && (
                                <>
                                    <span
                                        className={styles.accept}
                                        onClick={() => handleResponse("accept")}
                                    >
                                        Accept
                                    </span>
                                    <span
                                        className={styles.decline}
                                        onClick={() => handleResponse("decline")}
                                    >
                                        Decline
                                    </span>
                                </>
                            )}
                        </div>
                    </div>
                </div>
            )}
        </PopupContext.Provider>
    );
}
