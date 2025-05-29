import { useRouter } from "next/navigation"

export default function Header({ pageName, ele }) {
    const router = useRouter();
    return (
        <div className={`header df gp12 center ${ele ? "spB" : ""}`}>
            <div className="df gp6 center">
                <img className="icn1" src='/images/backIcon.png' onClick={() => { router.push('/home') }}></img>
                <h3>{pageName}</h3>
            </div>
            {ele}
        </div>
    )
}