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
                        className="avatarAndEmail"
                        onClick={() => router.push("/profile")}>

                        <Avatar url={user.avatar} name={user.firstName} />
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