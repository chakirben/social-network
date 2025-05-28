import React, { useEffect, useState } from "react"
import EventsList from "./eventList"
import "../../styles/eventInHome.css"
import Events from "./event"
import Divider from "../divider"


export default function HomeEvents() {
    const [events, setEvents] = useState([])

    useEffect(() => {
        const fetchEvents = async () => {
            const res = await fetch("http://localhost:8080/api/GetHomeEvents", {
                credentials: "include",
            })
            const data = await res.json()
            setEvents(data)
        }
        fetchEvents()
    }, [])

    return (
        <div className="EventInHome">
            <div className="type">Active Events</div>

            {events?.map((eve, index) => (
                <React.Fragment key={eve.id}>
                    <Events event={eve} index={index} />
                    <Divider />
                </React.Fragment>
            ))}
        </div>
    )

}
