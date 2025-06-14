import { useState, useEffect } from 'react';
import styles from './avatar.module.css';

export default function Avatar({ url, name, size }) {
    const [isValid, setIsValid] = useState(null);

    useEffect(() => {
        if (!url) {
            setIsValid(false);
            return;
        }

        const img = new Image();
        img.src = url;

        img.onload = () => setIsValid(true);
        img.onerror = () => setIsValid(false);
    }, [url]);

    const initial = name ? name[0].toUpperCase() : '?';
    const classNames = `${styles.letterAvatar} ${size || ''}`;

    if (isValid === true && url) {
        return <img className={`${styles.avatar} ${size || ''}`} src={url} alt={name || 'avatar'} />;
    }

    return (
        <div className={classNames}>
            <span>{initial}</span>
        </div>
    );
}
