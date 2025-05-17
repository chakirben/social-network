'use client'
import Header from "@/components/header";
import SideBar from "@/components/sidebar";
import "../home/home.css"
import { use, useState, useEffect } from "react";
import AcceptFollow from "./acceptFollow"


export default function NotifClient(session) {
    const [notifications, setNotifications] = useState();


    useEffect(() => {
        const fetchNotifications = async () => {
            try {
                const response = await fetch('http://localhost:8080/api/getNotifications', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ user_session: session.session }),
                });

                if (!response.ok) {
                    const errorText = await response.text();
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
                                <li key={index} className="notification">
                                    {new Date(notif.Date).toLocaleString()} : {notif.Sender} {notif.Type === "follow_request" && <span>Sends you a follow request</span>}
                                    {notif.Status === "pending" && notif.Type === "follow_request" && <AcceptFollow followerID = { notif.SenderID } followedSession = { session } />}
                                </li>
                            ))
                        )}
                    </ul>
                </div>
            </div>
        </div>
    )
}
