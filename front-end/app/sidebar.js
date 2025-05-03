import Link from 'next/link';
import './globals.css';

export default function SideBar() {
    return (
        <aside className="sidebar">
            <nav>
                <Link href="/home">Home</Link>
                <Link href="/groups">Groups</Link>
                <Link href="/chat">Chat</Link>
                <Link href="/users">Users</Link>
            </nav>
        </aside>
    )
}
