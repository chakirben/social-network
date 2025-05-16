import { useEffect, useState } from "react"
import "../styles/events.css"

export default function Events(groupId) {

    const [events, setEvents] = useState([])
    useEffect(() => {
        const fetchEvent = async () => {
            try {
                let res = await fetch(`http://localhost:8080/api/GetEvents?id=${groupId}`, {
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

    const handleIsGoing = async (eventId, isGoing, groupId) => {
        try {
            const res = await fetch("http://localhost:8080/api/SetAttendance", {
                method: 'POST',
                credentials: 'include',
                headers: {
                    "Content-Type": 'application/json'
                },
                body: JSON.stringify({
                    eventId: eventId,
                    isGoing: isGoing,
                    groupId: groupId
                })

            })
            console.log(eventId);

            if (res.ok) {
                setEvents((prev) =>
                    prev.map((eve) =>
                        eve.id === eventId ? { ...eve, isUserGoing: isGoing } : eve
                    ))
            } else {
                console.error("Failed to update attendance");
            }

        } catch (error) {
            console.error("-----errror-----", error);

        }

    }

    return (
        <>
            {events.map((event, index) => (
                <div className="event" key={index}>
                    <div className="eventData">
                        <p className="EventTitle">{event.title}</p>
                        <p className="Eventdescription">{event.description}</p>
                        <div className="dateAndName">
                            <div className="iconAndDate">
                                <img src="http://localhost:8080/uploads/dateOfBirth.svg" />
                                <p className="eventDate">{new Date(event.eventDate).toDateString()}</p>
                            </div>
                            <hr></hr>
                            <div className="avatarAndName">
                                <img className="eventAvatar" src={`http://localhost:8080/${event.avatar}`}></img>
                                <p className="eventCreator">{event.firstName} {event.firstName}</p>
                            </div>
                        </div>
                    </div>
                    <div className="twoBtns">
                        <span className={`goingBtn ${event.isUserGoing === true ? "selected" : ""}`}
                            onClick={() => handleIsGoing(event.id, true, event.groupId)}
                        >Going</span>
                        <span className={`notGoingBtn ${event.isUserGoing === false ? "selected" : ""}`}
                            onClick={() => handleIsGoing(event.id, false, event.groupId)}
                        >Not Going</span>
                    </div>
                </div>
            ))}
        </>
    )
}