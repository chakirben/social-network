import "../../styles/events.css"
import Divider from "../divider"
import Event from "./event"

export default function EventsList({ events }) {

    return (
        <>
            {events?.map((evnt, i) => (
                <Event event={evnt} key={i}></Event>
            ))}
            <Divider />
        </>
    )
}
