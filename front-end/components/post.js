export default function Post ({pst}) {
    return (
        <div className="post">
            <div className="postContent">
                <div className="userData">
                    <img className="pic sm" src={pst.image}></img>
                    <h4>{pst.creator}</h4>
                    <h5>{pst.created_at}</h5>
                </div>
                <div className="content">{pst.content}</div>
                <img className="pic nrml" src={pst.image}></img>
            </div>
        </div>
    )
}