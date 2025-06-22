import { timePassed } from "@/public/utils/timePassed";
import ReactionGroup from "./reactionGroup";
import Avatar from "./avatar/avatar";
import Divider from "./divider";

export default function Comment({ comment }) {

    const { id, content, image, firstName, lastName, avatar, createdAt } = comment;
    const fullName = `${firstName} ${lastName}`;

    return (
        <article className="comment">
            <div className="content">
                <div className="headerContainer">
                    <Avatar url={avatar} name={firstName} />
                    <div className="commentHeader">
                        <div className="authorAndTime">
                            <p className="commentAuthor">{fullName}</p>
                            <time className="commentTime" dateTime={createdAt}>
                                {"• " + timePassed(createdAt)}
                            </time>
                        </div>
                        <p className="commentContent">{content}</p>
                    </div>
                </div>
                {image && (
                    <img
                        className="commentImage"
                        src={image}
                        alt="Attached to comment"
                    />
                )}
            </div>
            <ReactionGroup className="reactionsContainer df gp12" likeCount={comment.likeCount} dislikeCount={comment.dislikeCount} itemType="comment" itemId={comment.id} userReaction={comment.userReaction} ></ReactionGroup>
            <Divider />
        </article>
    );
}
