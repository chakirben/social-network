import styles from './avatar.module.css';
export default function Avatar({name}) {
    console.log(name);
    
    return (
        <div className={styles.letterAvatar}>
            <span>{name[0]}</span>
        </div>
    )
}