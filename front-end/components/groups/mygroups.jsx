export default function MyGroup({ group , onView }) {
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
                <button id="buttontoseegroup" onClick={() => onView(group.Id)}>View</button>
            </div>
        </div>
    );
}

