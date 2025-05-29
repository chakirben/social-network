import { useRouter } from "next/navigation";
import { useUser } from "./userContext";
export default function ProfileButton() {
    const router = useRouter();
    const { user, setUser } = useUser();
    console.log(user);
    return (
        <div className="profileButton" onClick={() => { router.push("/profile") }}>
            {user ? (
                <div className="profileInfo df gp12 center">
                    <img className="avatar" src={"http://localhost:8080/" + user.avatar} />
                    <div className="profileDetails df col gp6">
                        <h4>{user.firstName}</h4>
                        <p>{user.email}</p>
                    </div>
                </div>
            )
                : ""}
        </div>
    );
}