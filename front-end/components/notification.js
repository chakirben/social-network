import "../styles/notification.css"
import { timePassed } from "@/public/utils/timePassed"

export default function Notification({ notification }) {
    let { id, firstName, notificationType, image, prjOrEvent, time } = notification
    const typeMsg = {
        "event": `created an event ${prjOrEvent}`,
        "goupInvite": `invited you to join ${prjOrEvent}`,
        "eventRequest": `requested access to ${prjOrEvent}`
    }

    const msg = typeMsg[notificationType] || "sent you a notification"
    const handleRes = async function (action) {
        try {
            let res = await fetch(`http://localhost:8080/api/Notification?id=${id}&action=${action}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ id })
            })
            if (!res.ok) {
                console.error(`Failed to ${action}:`, await response.text())
            } else {
                console.error(`Notification${action}ed seccesfully`)
            }
        } catch (error) {
            console.error('Error sending response:', error);
        }
    }

    return (
        <div className="notification">
            {/* {id} */}
            <div className="notifHeader">
                <img className="notifImage" src={`http://localhost:8080${image}`}></img>
                <div > <strong>{firstName}</strong> </div>
                <div className="notifMsg">{msg}</div>
                <time> • {timePassed(time)}</time>
            </div>
            <div className="buttons">
                <button onClick={() => handleRes("decline")}>Decline</button>
                <button onClick={() => handleRes("accept")}>Accept</button>
            </div>
        </div >
    )
}