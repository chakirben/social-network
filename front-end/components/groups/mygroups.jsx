import { useRouter } from 'next/navigation';
export default function MyGroup({ group }) {
    console.log(group.Id);
    const router = useRouter();

    return (
        <div className="groupc">
            <div className="groupContent">
                <h3>{group.Title}</h3>
                <p className="discriptiong">{group.Description}</p>
                <div className="usersandpost">
                    <img src="./images/users.svg"></img>
                    <p>{group.MembersCount}Users</p>
                    <img src="./images/postgroups.svg"></img>
                    <p> {group.PostCont} Posts</p>
                </div>
            </div>
            <div className="buttonjoin">
                <button className='secondaryBtn' href={`/groups/${group.Id}`} onClick={() => router.push(`/groups/${group.Id}?title=${group.Title}`) }>View</button>
            </div>
        </div>
    );
}

