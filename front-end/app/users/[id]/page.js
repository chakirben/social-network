import { cookies } from 'next/headers';
import ProfileClient from './profileClient';

export default async function ProfileServer({ searchParams }) {
    
    const cookieStore = cookies();
    const session = cookieStore.get('sessionId');

    return <ProfileClient session={session.value} searchParams={searchParams} />;
}