import "../../styles/eventInGroup.css"
import Avatar from "../avatar/avatar";
import Divider from "../divider"

export default function EventsList({ events, onRespond }) {
    if (!Array.isArray(events) || events.length === 0) {
        return (
            <div className="noEvents">
                <img className="noContent" src="/images/noContent.svg" alt="No content" />
                No events created, be the first
            </div>
        );
    }

    return (
        <>
            {events.map((event, index) => (
                <div key={index}>
                    <div className="event">
                        <div className="eventData">
                            <p className="EventTitle">{event.title}</p>
                            <p className="Eventdescription">{event.description}</p>
                            <div className="dateAndName">
                                <div className="iconAndDate">
                                    <img src="http://localhost:8080/uploads/dateOfBirth.svg" alt="date icon" />
                                    <p className="eventDate">{new Date(event.eventDate).toDateString()}</p>
                                </div>
                                <hr />
                                <div className="avatarAndName">
                                   <Avatar url={event.avatar} name={event.title} />
                                    <p className="eventCreator">{event.firstName} {event.lastName}</p>
                                </div>
                            </div>
                        </div>
                        <div className="twoBtns">
                            <span
                                className={`goingBtn ${event.isUserGoing === true ? "selected" : ""}`}
                                onClick={() => onRespond(event.id, true, event.groupId)}
                            >
                                Going
                            </span>
                            <span
                                className={`notGoingBtn ${event.isUserGoing === false ? "selected" : ""}`}
                                onClick={() => onRespond(event.id, false, event.groupId)}
                            >
                                Not Going
                            </span>
                        </div>
                    </div>
                    <Divider />
                </div>
            ))}
        </>
    );
}
