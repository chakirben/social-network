'use client';

import { useState } from "react";

export default function FollowButton({ follow_status,session, id }) {

  const  [btnText, setbtnText] = useState(follow_status);
  const handleFollow = async (e) => {
    e.preventDefault();
    const response = await fetch(`/api/follow?id=${id}&action=${btnText}`, {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        follower_session: session,
        followed_id: id,
      }),
    });

    if (!response.ok) {
      const errorText = await response.text();
      console.error('Error:', errorText);
    } else {
      const data = await response.text();
     setbtnText(data);
      
    }
  };

  return <button className="follow-button" onClick={handleFollow}>{btnText}</button>;
}
