import { useEffect, useState } from "react"
import EventsList from "./eventList"


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

    return (
        <div>
            <EventsList events={events} onRespond={handleRespond} />
        </div>
    )
}
