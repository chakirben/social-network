
import { cookies } from 'next/headers';
import FollowButton from './FollowButton';


export default async function Profile({ searchParams }) {
    const cookieStore = cookies();
    const session = cookieStore.get('session_id');
    //console.log(session.value);

    const name = searchParams?.name ?? 'anonymous';

  
    const response = await fetch('http://localhost:8080/api/profile', {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ nickname: "iafriad" }),
    });
  
    if (!response.ok) {
      const resp = await response.text();
      console.error('Error get Data:', resp);
    } else {
      const data = await response.json();
      console.log('Success:', data);
    }



    return (
      <div>
        <FollowButton session={session.value} name={name} />
      </div>
    );
}