import Avatar from "../avatar/avatar"

export default function Events({ event  , index}) {
    if(!event)return null
    return (
        <div className="event" key={index}>
            <div className="eventData">
                <p className="EventTitle">{event.title}</p>
                <p className="Eventdescription">{event.description}</p>
                <div className="dateAndName">
                    <div className="iconAndDate">
                        <img src="http://localhost:8080/uploads/dateOfBirth.svg" />
                        <p className="eventDate">{new Date(event.eventDate).toDateString()}</p>
                    </div>
                    <hr/>
                    <div className="avatarAndName">
                        <Avatar url={event.avatar} name={ event.title} size={"xs"}/>
                        <p className="eventCreator">{event.firstName} {event.lastName}</p>
                    </div>
                </div>
            </div>
        </div>
    )
}