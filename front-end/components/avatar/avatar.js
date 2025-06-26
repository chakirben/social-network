import { useState, useEffect } from 'react';
import styles from './avatar.module.css';

export default function Avatar({ url, name, size }) {
    console.log(url);
    
    const [isValid, setIsValid] = useState(null);
    useEffect(() => {
        if (!url) {
            setIsValid(false);
            return;
        }
        const img = new Image();
        img.src = url
        img.onload = () => setIsValid(true);
        img.onerror = () => setIsValid(false);
    }, [url]);

    const initial = name ? name[0].toUpperCase() : '?';

    let sizeClass = styles.md;
    if (size === 'xs') sizeClass = styles.xs;
    else if (size === 'big') sizeClass = styles.big;

    const letterClass = `${styles.letterAvatar} ${sizeClass}`;
    const avatarClass = `${styles.avatar} ${sizeClass}`;

    if (isValid === true && url) {
        return <img className={avatarClass} src={url} alt={name || 'avatar'} />;
    }

    return (
        <div className={letterClass}>
            <span>{initial}</span>
        </div>
    );
}
