import { useEffect, useState } from "react"
import { useParams } from "next/navigation"
import EventsList from "./eventList"

export default function GroupEventsPage({id, events, setEvents})  {
    
    const handleRespond = async (eventId, isGoing, groupId) => {
        const res = await fetch("http://localhost:8080/api/SetAttendance", {
            method: "POST",
            credentials: "include",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ eventId, isGoing, groupId })
        })

        if (res.ok) {
            setEvents(prev =>
                prev.map(e => e.id === eventId ? { ...e, isUserGoing: isGoing } : e)
            )
        }
    }

    return <EventsList events={events} onRespond={handleRespond} />
}