export default function Goups ({grp}) {
    return (
        <div className="group">
            <div className="groupContent">
                <h3>{grp.title}</h3>
                <p>{grp.description}</p>
                <div>
                    <i></i> <p>Users</p>
                    <i></i> <p>Posts</p>
                </div>  
            </div>
            <button>hiii</button>
        </div>
    )
}