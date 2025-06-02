import { useRouter } from "next/navigation"
import styles from './header.module.css';

export default function Header({ pageName, ele }) {
  const router = useRouter();
  return (
    <div className={styles.header}>
      <div className="df gp6">
        <img
          className={styles.icn1}
          src='/images/backIcon.png'
          alt="Back"
          onClick={() => router.push('/home')}
        />
        <h3>{pageName}</h3>
      </div>
      {ele}
    </div>
  )
}