import "../styles/notifPopUp.css"
export default function NotifPopUp({ senserUsername = 'userName', message = 'flan invite you to join (...) ', notifType }) {
    console.log("NotifPopUp rendered with type:", notifType);
    
    return (
        <div className="NotifPopUp">
            <div className="userName_avatar_message">
                <img className="SenderImage" src="../images/image.svg" />
                <div className="userName_message">
                    <p className="userName">{senserUsername}</p>
                    <p className="message">{message} </p>
                </div>
            </div>
            <div className="twoButtons">
                {notifType === "follow" || notifType === "groupInvite" ? (
                    <>
                        <span>Accept</span>
                        <span>Decline</span>
                    </>
                ) : (
                    <>
                        <span>Going</span>
                        <span>Not Going</span>
                    </>
                )}
            </div>
        </div >
    )
}