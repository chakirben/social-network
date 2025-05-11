'use client';

export default function FollowButton({ session, id }) {
  

  const handleFollow = async (e) => {
    e.preventDefault();

    console.log("this is the session :", session);
    console.log("this is the id of the follower :", id);

    const response = await fetch('http://localhost:8080/api/follow', {
      method: 'POST',
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
      const data = await response.json();
      console.log('Success:', data);
    }
  };

  return <button onClick={handleFollow}>Follow</button>;
}
