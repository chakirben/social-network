
export default function UserNotifications({ notification }) {

    return (
        <li className="notificationContainer">
            <div className="notifHeader">
                {notification.Avatar ? (
                    <img
                        className="userAvatar"
                        src={`http://localhost:8080/${notification.Avatar}`}
                    />
                ) : (
                    <div className="letterAvatar">
                        <span>{notification.Sender[0].toUpperCase()}{notification.Sender.split(" ")[1][0].toUpperCase()}</span>
                    </div>
                )}
                <span> <strong>{notification.Sender}</strong> </span>
                <span className="notifMsg">
                    {notification.Type === "follow_request" ? "sent you a follow request" : ""}
                </span>
                <span> <time> {notification.Date}</time> </span>
            </div>
        </li>
    )

}