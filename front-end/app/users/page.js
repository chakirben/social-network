
import { cookies } from 'next/headers';

export default async function Users() {
    const cookieStore = cookies();
    const session = cookieStore.get('session_id');

  
    const response = await fetch('http://localhost:8080/api/users', {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ session: session.value }),
    });
  
    if (!response.ok) {
        const resp = await response.text();
        console.error('Error get Data:', resp);
    } else {
        const data = await response.json();
        console.log('Success:', data.unfollowed_users);

        return (
            <div>
                <h1>User List</h1>
                <ul>
                    {data.unfollowed_users.map((user) => (
                    <li key={user.ID} style={{ cursor: 'pointer' }}>
                        <a href={'/profile/'+user.ID}> {user.Firstname} {user.Lastname} ({user.Nickname}) </a>
                    </li>
                    ))}
                </ul>
            </div>
        );
    }
}


