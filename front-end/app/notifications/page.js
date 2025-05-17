

import { cookies } from 'next/headers';
import NotifClient from './notifClient';

export default function Notifications() {
    const cookieStore = cookies();
    const session = cookieStore.get('sessionId');

    return <NotifClient session={session.value} />;
}