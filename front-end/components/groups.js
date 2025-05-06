export default function Group({ group }) {
    return (
           <div className="groupc">
                <div className="groupbar">
                    <i></i>Groups<p></p>
                </div>
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

// export default function groupbar() {
//     return (
//         <div className="groupbar">
//             <i></i>Groups<p></p>
//         </div>
//     )
// }
