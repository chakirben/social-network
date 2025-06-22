

export default function AcceptFollow({ followerID, followedSession }) {

    
    const acceptFollow = async (e) => {
        e.preventDefault();

        const response = await fetch(`/api/acceptFollowRequest`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                follower_id: followerID,
                followed_session: followedSession.session,
            }),
        });

        if (!response.ok) {
            const errorText = await response.text();
            console.error('Error:', errorText);
        } else {
            const data = await response.json();
            console.log('Success:', data);
            document.querySelector(".notif-buttons").style.display = "none";
        }
    }


    const declineFollow = async (e) => {
        e.preventDefault();

        const response = await fetch(`/api/declineFollowRequest`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                follower_id: followerID,
                followed_session: followedSession.session,
            }),
        });

        if (!response.ok) {
            const errorText = await response.text();
            console.error('Error:', errorText);
        } else {
            const data = await response.json();
            console.log('Success:', data);
            document.querySelector(".notif-buttons").style.display = "none";
        }
    }


    return (
        <div className="notif-buttons">
            <button onClick={acceptFollow}>Accept</button>
            <button onClick={declineFollow}>Decline</button>
        </div>
    )
}