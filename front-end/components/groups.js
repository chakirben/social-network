export default function Group({ group }) {
    return (
        <div className="groupc">
            <div className="groupContent">
                <h3>{group.Title}</h3>
                <p>{group.Description}</p>
                <div>
                    <i></i> <p>{group.MembersCount} Users</p>
                    <i></i> <p>Posts</p>
                </div>
            </div>
            <button>Join</button>
        </div>
    );
}

/*
 <div className="groupbar">
                   <img src="./images/arrow-left.svg" ></img>
                   Groups
                </div>
*/