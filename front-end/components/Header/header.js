import { useRouter, usePathname } from "next/navigation"
import styles from './header.module.css'
export default function Header({ pageName, ele }) {
  const router = useRouter();
  const pathname = usePathname();

  const handleBack = () => {
    console.log('paath', pathname);

    if (pathname.startsWith('/post')) {
      router.push('/home');
      return;
    }

    const parts = pathname.split('/').filter(Boolean);

    if (parts.length > 1) {
      const parentPath = '/' + parts.slice(0, parts.length - 1).join('/');
      router.push(parentPath);
    } else if (parts.length === 1) {
      router.push('/');
    } else {
      router.push('/home');
    }
  };

  return (
    <div className={styles.header + " center"}>
      <div className="df gp6">
        <img
          className={styles.icn1}
          src='/images/backIcon.png'
          alt="Back"
          onClick={handleBack}
        />
        <h3>{pageName}</h3>
      </div>
      {ele}
    </div>
  );
}
