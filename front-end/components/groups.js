export default function Groups ({grp}) {
    return (
        <div className="group">
            <div className="groupContent">
                <h3>{grp.Title}</h3>
                <p>{grp.Description}</p>
                <div>
                    <i></i> <p>{grp.MembersCount} Users</p>
                    <i></i> <p>Posts</p>
                </div>  
            </div>
            <button>hiii</button>
        </div>
    )
}