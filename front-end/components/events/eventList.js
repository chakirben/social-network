import "../../styles/events.css"
import Divider from "../divider"

export default function EventsList({ events, onRespond }) {
    return (
        <>
            {events?.map((event, index) => (
                <div className="event" key={index}>
                    <div className="eventData">
                        <p className="EventTitle">{event.title}</p>
                        <p className="Eventdescription">{event.description}</p>
                        <div className="dateAndName">
                            <div className="iconAndDate">
                                <img src="http://localhost:8080/uploads/dateOfBirth.svg" />
                                <p className="eventDate">{new Date(event.eventDate).toDateString()}</p>
                            </div>
                            <hr />
                            <div className="avatarAndName">
                                <img className="eventAvatar" src={`http://localhost:8080/${event.avatar}`} />
                                <p className="eventCreator">{event.firstName} {event.lastName}</p>
                            </div>
                        </div>
                    </div>
                    <div className="twoBtns">
                        <span className={`goingBtn ${event.isUserGoing === true ? "selected" : ""}`}
                            onClick={() => onRespond(event.id, true, event.groupId)}
                        >Going</span>
                        <span className={`notGoingBtn ${event.isUserGoing === false ? "selected" : ""}`}
                            onClick={() => onRespond(event.id, false, event.groupId)}
                        >Not Going</span>
                    </div>
                </div>
            ))}
            <Divider />
        </>
    )
}
