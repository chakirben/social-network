import { useRouter } from "next/navigation"

export default function Header({pageName}) {
    const router = useRouter();
    return (
        <div className="header df center">
            <img className="icn1" src='/images/backIcon.png' onClick={()=>{router.push('/home')}}></img>
            <h3>{pageName}</h3>
        </div>
    )
}