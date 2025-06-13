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

    if (isValid === null) {
        return <div className={styles.letterAvatar +` ${size}`}><span>{name[0]}</span></div>;
    }

    return isValid ? (
        <img className={styles.avatar+` ${size}`} src={url}/>
    ) : (
        <div className={styles.letterAvatar+` ${size}`}>
            <span>{name[0]}</span>
        </div>
    );
}
