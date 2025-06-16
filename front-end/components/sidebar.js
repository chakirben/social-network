'use client'
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useState } from 'react';
import '../styles/global.css';

export default function SideBar() {
    const pathname = usePathname();
    const [collapsed, setCollapsed] = useState(false);

    const toggleSidebar = () => {
        setCollapsed(!collapsed);
    };

    return (
        <aside className={`sidebar ${collapsed ? 'collapsed' : ''}`}>
            <div className='logoAndHide'>
                <img className='logo' src='/images/logo.svg' alt='Logo' />
                <img
                    src='/images/hideMenu.svg'
                    alt='Hide menu'
                    onClick={toggleSidebar}
                    className={`hideIcon ${collapsed ? 'rotated' : ''}`}
                />
            </div>
            <nav className="navBar">
                <Link href="/home" className={`navLink ${pathname === '/home' ? 'active' : ''}`}>
                    <img src='/images/home-icon.svg' alt='Home icon' />
                    {!collapsed && 'Home'}
                </Link>
                <Link href="/groups" className={`navLink ${pathname === '/groups' ? 'active' : ''}`}>
                    <img src='/images/groups-icon.svg' alt='Groups icon' />
                    {!collapsed && 'Groups'}
                </Link>
                <Link href="/chat" className={`navLink ${pathname === '/chat' ? 'active' : ''}`}>
                    <img src='/images/chat-icon.svg' alt='Chat icon' />
                    {!collapsed && 'Chat'}
                </Link>
                <Link href="/users" className={`navLink ${pathname === '/users' ? 'active' : ''}`}>
                    <img src='/images/users-icon.svg' alt='Users icon' />
                    {!collapsed && 'Users'}
                </Link>
                <Link href="/notifications" className={`navLink ${pathname === '/notifications' ? 'active' : ''}`}>
                    <img src='/images/notification-icon.png' alt='Users icon' />
                    {!collapsed && 'Notifications'}
                </Link>
            </nav>
        </aside>
    )
}