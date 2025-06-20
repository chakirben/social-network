'use client'
import Header from "@/components/Header/header";
import SideBar from "@/components/sidebar";
import "../home/home.css"
import { use, useState, useEffect } from "react";
import AcceptFollow from "./acceptFollow"
import UserNotifications from "@/components/Notification/notification";
import React from "react";


export default function NotifClient(session) {
    const [notifications, setNotifications] = useState();


    useEffect(() => {
        const fetchNotifications = async () => {
            try {
                const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/getNotifications`, {
                    method: 'POST',
                    credentials: 'include',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ user_session: session.session }),
                });

                if (!response.ok) {
                    const errorText = await response.json();
                    console.error('Error fetching data:', errorText);
                } else {
                    const data = await response.json();
                    setNotifications(data);
                    console.log("success data :", data);
                }
            } catch (error) {
                console.error('Network error:', error);
            }
        };
        fetchNotifications();
    }, []);

    if (!notifications) {
        return <div>Loading profile...</div>;
    }

    const { notif_data } = notifications;


    return (
        <div className="home">
            <SideBar />
            <div className="homeP">
                <Header pageName={'Notifications'} />
                <div className="usersList">
                    <ul className="notification-list">
                        {!notif_data ? (
                            <p>No notifications</p>
                        ) : (
                            notif_data.map((notif, index) => (
                                <React.Fragment key={index}>
                                    <ul>
                                    <UserNotifications notification={notif} />
                                    {notif.Status === "pending" && notif.Type === "follow_request" && (
                                        <AcceptFollow
                                            followerID={notif.SenderID}
                                            followedSession={session}
                                        />
                                    )}
                                    </ul>
                                </React.Fragment>
                            ))
                        )}
                    </ul>
                </div>
            </div>
        </div>
    )
}
