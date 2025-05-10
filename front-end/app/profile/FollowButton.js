'use client';

export default function FollowButton({ session, name }) {
  const handleFollow = async (e) => {
    e.preventDefault();

    console.log("this is the session :", session);
    console.log("this is the name of the follower :", name);

    const response = await fetch('http://localhost:8080/api/follow', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        follower_session: session,
        followed: name,
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
