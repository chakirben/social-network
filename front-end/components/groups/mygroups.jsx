import Link from "next/link";
export default function MyGroup({ group, onView }) {
    console.log(group.Id);

    return (
        <div className="groupc">
            <div className="groupContent">
                <h3>{group.Title}</h3>
                <p className="discriptiong">{group.Description}</p>
                <div className="usersandpost">
                    <img src="./images/users.svg"></img>
                    <p>{group.MembersCount}Users</p>
                    <img src="./images/postgroups.svg"></img>
                    <p> 2 Posts</p>
                </div>
            </div>
            <div className="buttonjoin">
                <Link href={`/groups/${group.Id}`} onClick={() => onView(group.Id)}>View</Link>
            </div>
        </div>
    );
}

