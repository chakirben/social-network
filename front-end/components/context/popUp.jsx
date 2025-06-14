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
    const [popup, setPopup] = useState(null);
    const router = useRouter();

    useEffect(() => {
        let timer;
        if (popup) {
            timer = setTimeout(() => {
                setPopup(null);
            }, 240000);
        }
        return () => clearTimeout(timer);
    }, [popup]);

    const showPopup = (notification) => {
        console.log(notification.data);

        setPopup(notification.data);
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
                        <span className={styles.user}>{name}</span> Requested To follow you
                    </p>
                );
            case "groupInvite":
                return (
                    <p>
                        <span className={styles.user}>{name}</span> {"Invited you to Join Group" + notif.title}
                    </p>
                );
            case "groupJoinRequest":
                return (
                    <p>
                        <span className={styles.user}>{name}</span> Requested to enter your
                        Group {notif.title}
                    </p>
                );
            default:
                return <p>New notification</p>;
        }
    };

    return (
        <PopupContext.Provider value={{ showPopup }}>
            {children}
            {popup && (
                <div className={styles.popupContainer} onClick={handleClick}>
                    <div className={styles.popup}>
                        <div className={styles.info}>
                            <Avatar url={popup.avatar} name={popup.firstName} size="small" />
                            <div className={styles.text}>{getMessage(popup)}</div>
                        </div>
                        <div className={styles.buttons}>
                            <span className={styles.accept}>Accept</span>
                            <span className={styles.decline}>Decline</span>
                        </div>
                    </div>
                </div>
            )}
        </PopupContext.Provider>
    );
}