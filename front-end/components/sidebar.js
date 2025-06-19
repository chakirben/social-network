'use client'
import Link from 'next/link';
import { usePathname, useRouter } from 'next/navigation';
import { useContext, useState, useEffect } from 'react';
import '../styles/global.css';
import { WebSocketContext } from './context/wsContext';

export default function SideBar() {
    const pathname = usePathname();
    const [collapsed, setCollapsed] = useState(false);
    const [isMobile, setIsMobile] = useState(false);

    const {
        notifCounter,
        setNotifCounter,
        messagesCounter,
        setMessagesCounter
    } = useContext(WebSocketContext);

    const router = useRouter();

    // Track window width and update isMobile
    useEffect(() => {
        function handleResize() {
            setIsMobile(window.innerWidth <= 800);
        }

        handleResize(); // Set initially

        window.addEventListener('resize', handleResize);
        return () => window.removeEventListener('resize', handleResize);
    }, []);

    // When isMobile changes, force collapsed to true
    useEffect(() => {
        if (isMobile) {
            setCollapsed(true);
        }
    }, [isMobile]);

    // Toggle only works when not mobile
    const toggleSidebar = () => {
        if (!isMobile) {
            setCollapsed(!collapsed);
        }
    };

    return (
        <aside className={`sidebar ${collapsed ? 'collapsed' : ''}`}>
            <div className='logoAndHide'>
                <img
                    className='logo'
                    src='/images/logo.svg'
                    alt='Logo'
                    onClick={() => {
                        router.push('/home');
                    }}
                />
                <img
                    src='/images/hideMenu.svg'
                    alt='Hide menu'
                    onClick={toggleSidebar}
                    className={`hideIcon ${collapsed ? 'rotated' : ''} ${isMobile ? 'disabled' : ''
                        }`}
                    style={{ cursor: isMobile ? 'not-allowed' : 'pointer' }}
                />
            </div>
            <nav className='navBar'>
                <Link href='/home' className={`navLink ${pathname === '/home' ? 'active' : ''}`}>
                    <img src='/images/home-icon.svg' alt='Home icon' />
                    {!collapsed && 'Home'}
                </Link>
                <Link href='/groups' className={`navLink ${pathname === '/groups' ? 'active' : ''}`}>
                    <img src='/images/groups-icon.svg' alt='Groups icon' />
                    {!collapsed && 'Groups'}
                </Link>
                <Link href='/chat' className={`navLink ${pathname === '/chat' ? 'active' : ''}`}>
                    <img src='/images/chat-icon.svg' alt='Chat icon' />
                    {!collapsed && 'Chat'}
                    {messagesCounter ? <span className='notificationBadge'>{messagesCounter}</span> : null}
                </Link>
                <Link href='/users' className={`navLink ${pathname === '/users' ? 'active' : ''}`}>
                    <img src='/images/users-icon.svg' alt='Users icon' />
                    {!collapsed && 'Users'}
                </Link>
                <Link
                    href='/notifications'
                    className={`navLink ${pathname === '/notifications' ? 'active' : ''}`}
                >
                    <img src='/images/Notification.svg' alt='Users icon' />
                    {!collapsed && 'Notifications'}
                    {notifCounter ? <span className='notificationBadge'>{notifCounter}</span> : null}
                </Link>

        
                    <Link
                        href='/profile'
                        className={`navLink ${pathname === '/profile' ? 'active' : ''}`}
                    >
                        <img src='/images/user.svg' alt='Users icon' />
                        {!collapsed && 'Profile'}
                    </Link>
            </nav>
        </aside>
    );
}