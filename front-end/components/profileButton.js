import { useRouter } from "next/navigation";
import { useUser } from "./context/userContext";
import Avatar from "./avatar/avatar";
export default function ProfileButton() {
    const router = useRouter();
    const { user, setUser } = useUser();
    console.log(user);
    return (
        <div className="profileButton">
            {user ? (
                <div className="profileInfo df gp12 center" >
                    <div
                        className="Notifications" onClick={(e) => {
                            e.stopPropagation()
                            router.push("/notifications")
                        }}
                    ><img src="./images/notifications.svg"></img>
                    </div>


                    <div
                        className="avatarAndEmail"
                        onClick={() => router.push("/profile")}>
                            
                        <img
                            className="avatar"
                            src={`http://localhost:8080/${user.avatar}`}
                            alt="Avatar"
                        />
                        <div className="profileDetails df col gp6">
                            <h4>{user.firstName}</h4>
                            <p>{user.email}</p>
                        </div>
                    </div>
                </div>
            )
                : ""}
        </div>
    );
}