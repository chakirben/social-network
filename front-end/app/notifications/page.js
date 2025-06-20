"use client";

import Header from "@/components/Header/header";
import Notification from "@/components/Notification/notification";
import SideBar from "@/components/sidebar";
import { useEffect, useState } from "react";
import "../home/home.css";
export default function Notifications() {
    const [notifications, setNotifications] = useState([]);

    useEffect(() => {
        const fetchNotifications = async () => {
            try {
                const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/getNotifications`, {
                    method: "POST",
                    credentials: "include",
                    headers: {
                        "Content-Type": "application/json",
                    },
                });

                if (!response.ok) {
                    const errorText = await response.json();
                    console.error("Error fetching data:", errorText);
                } else {
                    const data = await response.json();
                    setNotifications(data);
                    console.log("success data:", data);
                }
            } catch (error) {
                console.error("Network error:", error);
            }
        };

        fetchNotifications();
    }, []);
    return (
        <div className="home">
            <SideBar />
            <div className="homeP">
                <Header pageName={"Notifications"} />
                <div>
                    {notifications ? (
                        notifications.map((notif, index) => (
                            <Notification key={index} notification={notif} />
                        ))
                    ) : (
                        <p style={{ padding: "16px" }}>No notifications</p>)}
                </div>
            </div>
        </div>

    );
}
