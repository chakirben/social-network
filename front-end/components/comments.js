export default function Comment({ comment }) {
    const { id, content, image, firstName, lastName, avatar, createdAt } = comment;
    const fullName = `${firstName} ${lastName}`;

    return (
        <article className="comment">
            <div className="headerContainer">
                {avatar && (
                    <img
                        className="commentAvatar"
                        src={`http://localhost:8080${avatar}`}
                        alt={`${fullName}'s avatar`}
                    />
                )}
                <div className="commentHeader">
                    <div className="authorAndTime">
                        <p className="commentAuthor">{fullName}</p>
                        <time className="commentTime" dateTime={createdAt}>
                            {"• " + new Date(createdAt).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                        </time>
                    </div>
                    <p className="commentContent">{content}</p>
                </div>
            </div>
            {image && (
                <img
                    className="commentImage"
                    src={`http://localhost:8080${image}`}
                    alt="Attached to comment"
                />
            )}
        </article>
    )
}
