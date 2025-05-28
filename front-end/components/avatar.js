


export default function Avatar({user}) {


    return (
        <div className="letterAvatar">
            <span>{user.firstName[0].toUpperCase()}{user.lastName[0].toUpperCase()}</span>
        </div>
    )
}