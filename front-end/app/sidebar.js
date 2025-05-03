'use client'
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import './globals.css';

export default function SideBar() {
    const pathname = usePathname();
    console.log(pathname);


    return (
        <aside className="sidebar">
            <img className='logo' src='../images/logo.svg'></img>
            <nav className="navBar">
                <Link href="/home" className={`navLink ${pathname === '/home' ? 'active' : ''}`}>
                    <img src='../images/home-icon.svg'></img>
                    Home
                </Link>
                <Link href="/groups" className={`navLink ${pathname === '/groups' ? 'active' : ''}`}>
                    <img src='../images/groups-icon.svg'></img>
                    Groups
                </Link>
                <Link href="/chat" className={`navLink ${pathname === "/chat" ? 'active' : ''}`}>
                    <img src='../images/chat-icon.svg'></img>
                    Chat
                </Link>
                <Link href="/users" className={`navLink ${pathname === '/users' ? 'active' : ''}`}>
                    <img src='../images/users-icon.svg'></img>
                    Users
                </Link>
            </nav>
        </aside>
    )

}