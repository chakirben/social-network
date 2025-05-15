import { useEffect, useState } from "react"
import "../styles/events.css"


export default function Events() {

    const [events, setEvents] = useState([])

    useEffect(() => {
        const fetchEvent = async () => {
            try {
                let res = await fetch('http://localhost:8080/api/GetEvents', {
                    credentials: 'include',

                })
                const data = await res.json()
                setEvents(data)
            } catch (error) {
                console.error("Error fetching Events:", error);

            }
        }
        fetchEvent()
    }, [])
    return (
        <>
            {events.map((event, index) => (
                <div className="event" key={index}>
                    <div className="eventData">
                        <p className="EventTitle">{event.title}</p>
                        <p className="Eventdescription">{event.description}</p>
                        <div className="dateAndName">
                            <p className="eventDate">{event.eventDate}</p>
                            <p className="eventCreator">{event.firstName} {event.firstName}</p>
                        </div>
                    </div>
                    <div className="twoBtns">
                        <button>Going</button>
                        <button>Not Going</button>
                    </div>
                </div>
            ))}
        </>
    )
}